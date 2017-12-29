#!/bin/bash

go install github.com/happyyi008/remember

if [ -f $GOBIN/remember ]; then
    sudo mv $GOBIN/remember /usr/local/bin/rmb
else
    sudo mv $GOPATH/bin/remember /usr/local/bin/rmb
fi
