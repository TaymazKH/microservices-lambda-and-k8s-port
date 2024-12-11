#!/bin/bash -eu

protodir=../../protos
protoname=genproto

mkdir -p $protoname

protoc --js_out=import_style=commonjs,binary:./$protoname \
       -I $protodir \
       $protodir/demo.proto
