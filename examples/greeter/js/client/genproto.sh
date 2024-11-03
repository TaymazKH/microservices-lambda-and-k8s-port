#!/bin/bash -eu

protodir=../../protos
protoname=genproto

mkdir -p $protoname

protoc --js_out=import_style=commonjs,binary:./$protoname \
       -I $protodir \
       $protodir/hello.proto
#       --grpc_out=grpc_js:./$protoname \
#       --plugin=protoc-gen-grpc=$(which grpc_tools_node_protoc_plugin) \
# todo: check if the above comments are needed
