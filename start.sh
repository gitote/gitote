#!/bin/bash

export GOPATH=$HOME/go && export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin && cd $GOPATH/src/gitote.com/gitote && cd gitote
git pull
go build
mysql-ctl start
./gitote web -p 8080