# gt

Simple git repo clone tool with workspace support

## Install

```sh
go install github.com/hnlq715/gt@latest
```

## Create `.gtconfig.yaml` in $HOME dir

Specify gitconfig for each host, which means you can use different email for different git host.

```yaml
gitconfig:
  - host: github.com
    email: hnlq.sysu@gmail.com
  - host: gitlab.com
    email: xyz@gitlab.com
```

## Usage

```sh
NAME:
   gt - Simple git repo clone tool with workspace support

USAGE:
   gt [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --workspace value  workspace path (default: "/home/pilot/workspace")
   --help, -h         show help
```
