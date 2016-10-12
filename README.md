Playment's Angel
================

Copyright 2016 Playment Inc.

WorkFlow MicroService written in [go](https://golang.org/)

This was planned to be a microservice for main playment server (crowdy) but
it turned out to be the parent of all.

# Installation

- Install & setup go from its [website](https://golang.org/)
- Create a workspace directory & clone the repo

```
$ mkdir -p $HOME/code/gocode/src/github.com/crowdflux
$ cd $HOME/code/gocode/src/github.com/crowdflux
$ git clone https://github.com/crowdflux/angel.git
```

- Setup GOPATH:

```
$ echo "export GOPATH='$HOME/code/gocode/'" >> ~/.bash_profile
$ echo "export GOENV='development'" >> ~/.bash_profile
$ echo "export PATH=$PATH:$GOPATH/bin" >> ~/.bash_profile
$ source ~/.bash_profile
```

- Install application:

```
$ cd $HOME/code/gocode/src/github.com/crowdflux/angel
$ git checkout master
$ go get -v
$ go get -t -v
```

- Install mongodb:

```
$ brew update
$ brew install mongodb
$ brew services start mongodb
```

- Create Index (mongo db name : playment_mongo_local)

```
$ mongo
> use playment_mongo_local
> db.feedline_input.createIndex({
  project_id : 1,
  reference_id : 1
  },
  {
          unique:true
  }
);
```



- Copy config file

```
$ cp app/config/development.example.json app/config/development.json
```

After copying edit the configuration file according to your local system




# Usage

To start the server:

```
$ go install
$ $GOPATH/bin/angel
```
