📦 protopkg
---

A declarative protocol buffer package manager. protopkg helps synchronize protocol buffers across multiple repositories. A `protopkg.json` file will be read in the current working directory and protocol buffer directories (or single files) will be pulled from github and copied to the desired path.

## Installation
Homebrew: `brew tap reverbdotcom/reverb && brew install protopkg`

With Go 1.11: `go get -u github.com/reverbdotcom/protopkg`

## Private Repositories
`protopkg` will make authorized calls to the public github API if you set a token via the command line. Ensure that this token has access to read the given repositories referenced in your `protopkg.json`

https://blog.github.com/2013-05-16-personal-api-tokens/

## Usage
```
NAME:
   protopkg - package manager for protocol buffers

USAGE:
   protopkg [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     token, t  set your GitHub token for pulling from private repositories
     local, l  sync a dependency based on the configured local path
     sync, s   pull down the protos - protopkg sync
     init, i   creates a new protopkg.json in the current directory - protopkg init
     add, a    adds a new proto dependency - protopkg add google/protos@HEAD ./protos/google
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## protopkg.json
```
{
  "protos": {
    "google/transit/gtfs-realtime/proto": {
      "path": "protos/gtfs",
      "ref": "db3c001d0c2412825c6911628ded36c583e60f95",
      "local": "../a-local-path"
    }
  }
}
```

