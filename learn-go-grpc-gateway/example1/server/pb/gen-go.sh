#!/usr/bin/env bash
set -eu

cd "$(dirname "$0")"

find . \( -type f  -name '*.pb.go' -o -name '*.pb.gw.go' -o -name '*.json' \) -exec rm -f {} +

export GO111MODULE="on"


go get \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc


GW_ORIGIN_VERSION=v1.14.4
go get github.com/grpc-ecosystem/grpc-gateway@$GW_ORIGIN_VERSION

GOPATH=$(go env GOPATH)
ARGS="\
--proto_path=. \
--proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$GW_ORIGIN_VERSION \
--proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$GW_ORIGIN_VERSION/third_party/googleapis \
--go_out=. \
--go_opt=paths=source_relative \
--go-grpc_out=. \
--go-grpc_opt=paths=source_relative \
--grpc-gateway_opt=logtostderr=true \
--grpc-gateway_out=. \
--grpc-gateway_opt=paths=source_relative \
--grpc-gateway_opt=repeated_path_param_separator=ssv"

for dir in $(find . -name '*.proto' | xargs -I{} dirname {} | sort | uniq); do
    echo "building $dir/*.proto"
    protoc $ARGS $dir/*.proto
done