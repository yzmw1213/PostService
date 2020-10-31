#!/bin/sh
DIR=$(cd "$(dirname "$0")" || exit 1 ; pwd);
POST_OUT_DIR="/grpc/postservice";
TAG_OUT_DIR="/grpc/tagservice";
POST_PROTO_FILE="post.proto";
TAG_PROTO_FILE="tag.proto";

protoc \
  --go_out=plugins=grpc:"${DIR}${POST_OUT_DIR}" \
  -I".${POST_OUT_DIR}" "${POST_PROTO_FILE}"

protoc \
  --go_out=plugins=grpc:"${DIR}${TAG_OUT_DIR}" \
  -I".${TAG_OUT_DIR}" "${TAG_PROTO_FILE}"
