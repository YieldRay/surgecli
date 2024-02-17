# surgecli

Third party [surge.sh](https://surge.sh) cli written in golang

Features:

-   Single executable
-   Compatible with the official CLI
-   Multi-user support
-   Friendly for CI environments

> The only difference is that we need use `surgecli upload` or `surgecli deploy` to deploy sites

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
# if you have already logged-in, run this
surgecli token
# or just fetch the token, but do not perform login
surgecli fetch-token <username> <password>
# for another machine, set environment variable
export SURGE_TOKEN=<your_token>
# use other command that require token without login
surgecli upload --silent ./dist mysite.surge.sh
```

Command help

```sh
NAME:
   surgecli - third party surge.sh cli

USAGE:
   surgecli [global options] command [command options]

COMMANDS:
   account, plan         Show account information
   api                   Set or Show API base URL, the official one is https://surge.surge.sh
   fetch-token           Fetch token by email and password, but do not save the token like login command
   list, ls              List my sites
   login                 Login (or create new account) to surge.sh
   logout                Logout from surge.sh
   su                    Switch user
   teardown, delete, rm  Delete site from current account
   token                 Show current token
   upload, deploy        Upload a directory (i.e. deploy a project) to surge.sh
   whoami                Show who you are logged in as
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
