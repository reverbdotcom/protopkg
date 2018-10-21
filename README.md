📦 protopkg
---

A declarative protocol buffer package manager. protopkg helps synchronize protocol buffers across multiple repositories. A `protopkg.json` file will be read in the current working directory and protocol buffer directories (or single files) will be pulled from github and copied to the desired path.

## Installation
Homebrew: `brew tap ebenoist/ebenoist && brew install protopkg`

With Go 1.11: `go get -u github.com/ebenoist/protopkg`

## Usage
```
NAME:
   protopkg - package manager for protocol buffers

USAGE:
   protopkg [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     sync, s  pull down the protos - protopkg sync
     init, i  creates a new protopkg.json in the current directory - protopkg init
     add, a   adds a new proto dependency - protopkg add google/protos@HEAD ./protos/google
     help, h  Shows a list of commands or help for one command

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
      "ref": "db3c001d0c2412825c6911628ded36c583e60f95"
    }
  }
}
```

## Private Repositories
`protopkg` will make authorized calls to the public github API if the environment variable `GITHUB_TOKEN` is present. Ensure that this token has access to read the given repositories referenced in your `protopkg.json`
