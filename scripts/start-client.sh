#!/bin/bash

# Start hello service
export GRPC_XDS_BOOTSTRAP=$PWD/example/bootstrap.json
go run ./example/cmd/client --name World \
    --addresses xds:///hello,xds:///world