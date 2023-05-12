# surgecli

third party [surge.sh](https://surge.sh) cli written in golang

```sh
USAGE:
   surgecli [global options] command [command options] [arguments...]

COMMANDS:
   login           login (or create new account) to surge.sh
   logout          logout to surge.sh
   account         show account information
   whoami          show my email
   list            list my sites
   upload, deploy  upload directory to surge.sh with specified domain
   teardown        delete site from surge.sh
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --API_HOST value  configure the API host (default: "https://surge.surge.sh")
   --help, -h        show help
```
