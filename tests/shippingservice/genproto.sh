#!/bin/bash -eu

protodir=../../protos
protoname=genproto

mkdir -p $protoname

protoc --go_opt=Mdemo.proto="/$protoname" \
       --go_opt=paths=source_relative \
       --go_out=./$protoname \
       -I $protodir \
       $protodir/demo.proto
