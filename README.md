# cli
Terminal client for Server. Quick visit https://binacs.cn

## 1. Preparation

Golang enviroment.

## 2. Usage

1.  Deploy `clid` which is the daemon program connected to *api.binacs.cn* , and `cli` which you use.

    **Run deploy.sh:**

    ```sh
    $ ./deploy.sh
    ```

2.  Quick visit binacs.cn serivice by the command line tool `cli`.

    **Run `cli --help` to see more details:**

    ```sh
    $ cli --help
    Terminal client for https://binacs.cn
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

## 3. More

1.  `cos` : Storage service, web at https://binacs.cn/toys/storage .
2.  `crypto` : Crypto service, web at https://binacs.cn/toys/crypto
3.  `pastebin` : PasteBin service, web at https://binacs.cn/toys/pastebin .
4.  `tinyurl` : TinyURL service, web at https://binacs.cn/toys/tinyurl .
5.  `user` : User service, only for cli.
6.  `version`: Show version infomation.
7.  `help` : You know.