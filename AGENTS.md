# AGENTS.md

Instructions for AI coding agents invoking the `cli` tool itself (not just
editing this repo's code) — e.g. to actually run `cos put` on behalf of a
user.

## Prefer the ephemeral Docker container

Don't install `clid` as a persistent daemon (`./deploy.sh`) to run a single
command on a user's behalf. Use the container described in README.md §2.1
instead — it leaves nothing behind: no binaries on `PATH`, no daemon, no
socket file, no pass key written anywhere.

The image is published to Docker Hub as `binacslee/cli` on every push to
`main` (`.github/workflows/docker.yml`) — pull it, don't build it, unless
you're specifically testing an unpublished local change:

```sh
docker run --rm -i \
  -v "<absolute-host-path>:/data/<name>:ro" \
  binacslee/cli <subcommand> <args...>
```

- Mount only the one file the command actually needs, read-only (`:ro`).
- Use `-i`, not `-t`/`-it` — there is no real terminal in an agent's shell
  tool. `cli` detects the non-TTY stdin automatically and reads one plain
  line instead of doing a hidden-echo prompt, so plain stdin redirection
  works correctly.

## Supplying the COS pass key

`cos put` / `cos get` require a shared secret (`CosConfig.PassKey` on the
server). Never hardcode it in a command, a file, or this document, and
never print it back once you have it. Expect the human operator to have
already exported it in their shell (e.g. `$COS_PASSKEY`) and pipe it
straight through:

```sh
docker run --rm -i \
  -v "$PWD/report.pdf:/data/report.pdf:ro" \
  binacslee/cli cos put /data/report.pdf <<< "$COS_PASSKEY"
```

If no such variable is set or referenced, stop and ask the human to
provide the pass key out-of-band — do not guess it, do not fall back to an
empty or placeholder value, and do not ask them to paste the raw secret
into chat if a shell env var will do.

## Before running `cos put`

This uploads to a real, publicly-readable production bucket
(`api.binacs.space`) — visible to anyone with the URL, billed to the
account owner, and not reversible through this tool (no delete command
exists client-side). Treat it as a hard-to-reverse, externally-visible
action:

- Confirm the exact file and intent with the user before running it,
  unless they've already explicitly authorized this specific upload
  earlier in the same conversation.
- Never invoke it speculatively, in a retry loop, or on a file whose
  contents you haven't confirmed are what the user intends to share
  publicly.

## Other subcommands

`crypto`, `pastebin`, `tinyurl`, `user`, `version` follow the same
container invocation pattern and don't need a pass key — only `cos` does.
