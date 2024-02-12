# surgecli

Third party [surge.sh](https://surge.sh) cli written in golang

Features:

-   Single executable
-   Compatible with the official CLI
-   Multi-user support
-   Friendly for CI environments

For github actions, see [deploy-to-surge-action](https://github.com/YieldRay/deploy-to-surge-action)

## usage

Manage your site

```sh
surgecli login <username> <password> # skip this if you have already logged in

# deploy one site
surgecli upload ./dist mydomain.example.net

# list my sites
surgecli list

# delete one site
surgecli teardown mydomain.example.net
```

You may want to upload your site with something like Github Actions, see this

```sh
# first, fetch token from your local machine
surgecli fetch-token <username> <password>
# if you have already logged-in, prefer this
surgecli fetch-token --local
# for another machine, set environment variable
export SURGE_TOKEN=<your_token>
# use other command that require token without login
surgecli upload --silent . mysite.surge.sh
```

Command help

```sh
NAME:
   surgecli - thrid party surge.sh cli

USAGE:
   surgecli [global options] command [command options]

COMMANDS:
   account, me           Show account information
   api                   Set or Show API base URL, the official one is https://surge.surge.sh
   fetch-token, token    Fetch token by email and password, but do not save the token like login command
   list, ls              List my sites
   login                 Login (or create new account) to surge.sh
   logout                Logout from surge.sh
   su                    Switch user
   teardown, delete, rm  Delete site from current account
   upload, deploy        Upload a directory (a.k.a. deploy a project) to surge.sh
   whoami                Show my email
   help, h               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug     toggle debug on (default: false)
   --help, -h  show help
```

## build

```sh
git clone https://github.com/YieldRay/surgecli.git
cd surgecli
go build -ldflags="-s -w" surgecli.go

# cross platform
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/surgecli-linux-amd64 -ldflags="-s -w" surgecli.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/surgecli-darwin-amd64 -ldflags="-s -w" surgecli.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/surgecli-windows-amd64.exe -ldflags="-s -w" surgecli.go
```
