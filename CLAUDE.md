# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is `binacs/cli`, the terminal client for the binacs.space service (https://binacs.space). It produces two binaries with a daemon/client split:

- **`clid`** — a background daemon that dials the remote gRPC API server (`api.binacs.space`) over TLS and re-exposes the same gRPC services locally over a Unix domain socket.
- **`cli`** — the user-facing command-line tool. It connects to `clid` via the local Unix socket (insecure local gRPC) and issues requests; it never talks to the remote API directly.

The actual protobuf service definitions and message types live in the external module `github.com/binacs/server` (imported as `cos_pb`, `crypto_pb`, `pastebin_pb`, `tinyurl_pb`, `user_pb`, plus `github.com/binacs/server/types`), not in this repo.

## Build / Deploy Commands

```sh
make            # clean + build; produces bin/clid and bin/cli
make build      # build only
make clean      # rm -rf bin
./deploy.sh     # installs cli/clid to ~/.local/bin (downloads the matching
                # prebuilt release binary, falling back to a local `go build`
                # only if none is reachable) and registers clid as a per-user
                # background service (LaunchAgent on macOS, `systemctl --user`
                # on Linux). No sudo.
docker build -t cli .   # build the ephemeral-runner image, see AGENTS.md
```

There are no automated tests in this repository (`go build ./...` / `go vet ./...` are the available checks).

Version metadata (`version.Maj`/`Min`/`Fix` in `version/version.go`, and `version.GitCommit`) is stamped at build time via `-ldflags` in the Makefile — do not hardcode `GitCommit`, it is injected by `git rev-parse HEAD`. `.github/workflows/release.yml` cross-compiles both binaries for darwin/linux × amd64/arm64 and publishes them to GitHub Releases on any `v*` tag push; that's what `deploy.sh` downloads from.

For any one-off `cli` invocation (yours or an agent's) that shouldn't install anything persistent, prefer the ephemeral Docker container over `deploy.sh` — see **AGENTS.md** for the exact command form and the pass-key handling rule.

## Architecture

### Daemon/client split over a Unix socket

`clid start` (cmd/clid/command/start.go):
1. Dials `api.binacs.space:30000` (flags: `--domain`, `--port`, `--instance`) using TLS with an embedded Cloudflare origin cert (`util.GetCertPool()` in util/crt.go). No credential is configured at Dial time; a client-side interceptor (`defaultAuthInterceptor`) stamps a non-secret placeholder `authorization` header onto any outbound call that doesn't already carry one — the server never validates it for non-Cos services (see server/gateway/grpc.go `auth()`, which extracts but never checks the token). Cos calls carry a real per-call secret instead, see below.
2. Listens on a Unix socket at `util.GetSockPath()` (util/sock.go — `$HOME/cli.sock`, falling back to `/var/run/cli.sock`).
3. Wraps that remote connection in `service.InitService(conn)` and re-registers each service (Cos, Crypto, Pastebin, TinyURL, User) as a local gRPC server over the socket — effectively proxying remote RPCs to local Unix-socket RPCs.

`cli` (cmd/cli/command/root.go), on every invocation, dials that same Unix socket insecurely (`unixConnect` + `grpc.WithInsecure()`) in `RootCmd.PersistentPreRunE`, then calls `service.InitService(conn)` again to get client stubs bound to the local socket. Subcommands then just call methods on the package-level `node *service.NodeServiceImpl`.

So `clid` must be running (see deploy.sh's LaunchAgent/systemd-user setup, or the Docker container) before `cli` commands will work.

### Cos pass-key relay (no secret ever touches disk)

`server`'s `CosServiceImpl.CosPut`/`CosGet` (server repo, service/cos.go) require the bearer token to equal `CosConfig.PassKey` — the same shared secret the web upload form (`gateway/web.go apiCosPut`) already enforced; other services are still unauthenticated in practice. On the `cli` side this secret is never written to disk or baked into a config file:
1. `cmd/cli/command/cos.go` calls `util.PromptSecret` to read it from the controlling terminal (echo off) right before a `cos put`, and attaches it to that one request via `util.AttachAuth` (util/auth.go).
2. `clid`'s `CosClientImpl.CosPut`/`CosGet` (service/cos.go) call `util.RelayAuth` to copy that header from the incoming local-socket request onto the outgoing call to the real server — `clid` never persists it either.

When adding new Cos-touching code paths, follow this relay pattern rather than reintroducing a stored credential.

### Dependency injection via `binacsgo/inject`

`service.InitService` (service/node.go) wires everything together using reflection-based DI (`github.com/binacsgo/inject`): each service impl is registered by name (`inject.Regist("Cos", ...)`) and struct fields tagged `inject-name:"X"` are populated by `inject.DoInject()`. This same wiring function is used by both `clid` (to get server implementations backed by the remote conn) and `cli` (to get client stubs backed by the local socket conn) — the only difference is which `*grpc.ClientConn` is passed in.

Each service in service/ (cos.go, crypto.go, pastebin.go, tinyurl.go, user.go) follows the same pattern: a `*Impl` struct with an injected `Conn *grpc.ClientConn`, an `AfterInject()` that constructs the pb client (`cos_pb.NewCosClient(impl.Conn)`), and methods that just forward to that pb client. service/interface.go declares the corresponding `XClient` interfaces consumed by `NodeServiceImpl`.

### Adding a new subcommand

Follow the existing pattern in cmd/cli/command/ (e.g. cos.go, tinyurl.go): one file per service, a single `cobra.Command` with an `Args` validator using `checkArgs`/`checkArgLength` (common.go), a `Run` that switches on a lowercased first-arg "op", and calls into `node.<Service>.<Method>(context.Background(), &pb.Req{...})` piped through `handleResp`. Register the new command in cmd/cli/main.go's `rootCmd.AddCommand(...)` list. If it's a new service (not just a new op on an existing one), add the corresponding client interface in service/interface.go, a `*ClientImpl` in service/, and wire it into `NodeServiceImpl`/`InitService` in service/node.go — mirroring what cmd/clid/command/start.go does to register the server side.
