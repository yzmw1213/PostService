#!/bin/sh
DIR=$(cd "$(dirname "$0")" || exit 1 ; pwd);
OUT_DIR="/grpc/post_grpc";
POST_PROTO_FILE="post.proto";
TAG_PROTO_FILE="tag.proto";

protoc \
  --go_out=plugins=grpc:"${DIR}${OUT_DIR}" \
  -I".${OUT_DIR}" "${POST_PROTO_FILE}"

protoc \
  --go_out=plugins=grpc:"${DIR}${OUT_DIR}" \
  -I".${OUT_DIR}" "${TAG_PROTO_FILE}"
