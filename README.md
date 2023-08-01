# surgecli

third party [surge.sh](https://surge.sh) cli written in golang  
currently, `.surgeignore` file is not fully supported

## usage

manage your site

```sh
surgecli login <username> <password> # skip this if you have already logged in

# deploy one site
surgecli upload ./dist mydomain.example.net

# list my sites
surgecli list

# delete one site
surgecli teardown mydomain.example.net
```

you may want to upload your site with something like Github Actions, see this

```sh
# first, fetch token from your local machine
surgecli fetch-token <username> <password>
# for another machine, set environment variable
export SURGECLI_TOKEN=<your_token>
# use other command that require token without login
surgecli upload --silent . mysite.surge.sh
```

command help

```
NAME:
   surgecli - thrid party surge.sh cli

USAGE:
   surgecli [global options] command [command options] [arguments...]

COMMANDS:
   account           Show account information
   fetch-token       Fetch token by email and password, but do not save the token like login command
   list              List my sites
   login             Login (or create new account) to surge.sh
   logout            Logout from surge.sh
   api               Print or set api host, the official host is https://surge.surge.sh
   su                Switch user
   teardown, delete  Delete site from surge.sh
   upload, deploy    Upload a directory (a.k.a. deploy a project) to surge.sh
   whoami            Show my email
   help, h           Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug     toggle debug on (default: false)
   --help, -h  show help
```

## build

```sh
git clone https://github.com/YieldRay/surgecli.git
cd surgecli
go build -ldflags="-s -w" surgecli.go

# build for linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" surgecli.go
```
