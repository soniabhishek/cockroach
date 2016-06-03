Playment's Angel Server
=======================

Copyright 2016 Playment Inc.

This is a MicroService written in [go](https://golang.org/)

# Installation

- Install & setup go from its [website](https://golang.org/)
- Create a workspace directory & clone the repo

```
$ mkdir -p $HOME/code/gocode/src/gitlab.com/playment-main
$ cd $HOME/code/gocode/src/gitlab.com/playment-main
$ git clone git@gitlab.com:playment-main/angel.git
```

- Setup GOPATH:

```
$ echo "export GOPATH='$HOME/code/gocode/'" >> ~/.bash_profile
$ echo "export PATH=$PATH:$GOPATH/bin" >> ~/.bash_profile
$ source ~/.bash_profile
```

- Install application:

```
$ cd $HOME/code/gocode/src/gitlab.com/playment-main/angel
$ git checkout develop
$ go get
$ go get -t
```

- Install mongodb:

```
$ brew update
$ brew install mongodb
$ brew services start mongodb
```

# Usage

It exposes Web Api which can be accessed from within the aws private server group.
Api Definitions can be found on server/api/builder.go

To start the server:

```
$ go install
$ $GOPATH/bin/support
```
