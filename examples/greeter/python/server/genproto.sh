#!/bin/bash -eu

protodir=../../protos
protoname=genproto

mkdir -p $protoname
touch $protoname/__init__.py

python -m grpc_tools.protoc \
       --python_out=./$protoname \
       --grpc_python_out=./$protoname \
       -I$protodir \
       $protodir/hello.proto