#!/bin/bash -eu

protodir=../../protos
protoname=genproto

mkdir -p $protoname

protoc --go_opt=Mgreeter.proto="/$protoname" \
       --go_opt=paths=source_relative \
       --go_out=./$protoname \
       --go-grpc_opt=Mgreeter.proto="/$protoname" \
       --go-grpc_opt=paths=source_relative \
       --go-grpc_out=./$protoname \
       -I $protodir \
       $protodir/greeter.proto
