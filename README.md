# cli
Terminal client for Server. Quick visit https://binacs.space

## 1. Preparation

None required. `deploy.sh` downloads a prebuilt binary for your OS/arch from the
[latest release](https://github.com/binacs/cli/releases/latest); a local Golang
environment is only needed as a fallback if no matching release asset is found.

## 2. Usage

1.  Deploy `clid` which is the daemon program connected to *api.binacs.space* , and `cli` which you use.

    **Run deploy.sh** (installs to `~/.local/bin`, registers `clid` as a
    per-user background service — no sudo required):

    ```sh
    $ ./deploy.sh
    ```

2.  Quick visit binacs.space serivice by the command line tool `cli`.

    **Run `cli --help` to see more details:**

    ```sh
    $ cli --help
    Terminal client for https://binacs.space
    More at https://github.com/binacs/cli
    
    Usage:
      root [command]
    
    Available Commands:
      cos         Cos Command:	 Just run `cli cos put/get sth.(file)`
      crypto      Crypto Command:	 Just run `cli crypto encrypt/decrypt BASE64/AES/DES sth.(string)`
      help        Help about any command
      pastebin    PasteBin Command:	 Just run `cli pastebin submit sth.(file)`
      tinyurl     TinyURL Command:	 Just run `cli tinyurl encode/decode sth.`
      user        User Command:	 Just run `cli user test/register/auth/refresh/info`
      version     Version Command
    
    Flags:
      -h, --help   help for root
    
    Use "root [command] --help" for more information about a command.
    ```

## 2.1 One-off use via Docker

For a single command with nothing left behind on the host afterwards — no
installed binaries, no `clid` daemon, no socket file, no typed pass key
saved anywhere — run it in an ephemeral container instead of `deploy.sh`.
The image is built and pushed to Docker Hub automatically on every push to
`main` (see `.github/workflows/docker.yml`), so there's normally nothing to
build locally:

```sh
$ docker run --rm -it \
    -v "$(pwd)/file.txt:/data/file.txt:ro" \
    binacslee/cli cos put /data/file.txt
```

The container starts `clid` in the background, runs the given `cli`
subcommand against it, and exits — `--rm` removes the container (and
everything `clid` created inside it) immediately after. `-it` is required
so the COS pass key prompt has a terminal to read from without echoing it;
only mount the specific file(s) the command needs.

To build locally instead (e.g. to test an unpublished change):

```sh
$ docker build -t cli .
$ docker run --rm -it -v "$(pwd)/file.txt:/data/file.txt:ro" cli cos put /data/file.txt
```

## 3. More

1.  `cos` : Storage service, web at https://binacs.space/toys/storage .
2.  `crypto` : Crypto service, web at https://binacs.space/toys/crypto
3.  `pastebin` : PasteBin service, web at https://binacs.space/toys/pastebin .
4.  `tinyurl` : TinyURL service, web at https://binacs.space/toys/tinyurl .
5.  `user` : User service, only for cli.
6.  `version`: Show version infomation.
7.  `help` : You know.