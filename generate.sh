#!/bin/sh
DIR=$(cd "$(dirname "$0")" || exit 1 ; pwd);
POST_OUT_DIR="/grpc/postservice";
TAG_OUT_DIR="/grpc/tagservice";
USER_OUT_DIR="/grpc/userservice";
POST_PROTO_FILE="post.proto";
TAG_PROTO_FILE="tag.proto";
USER_PROTO_FILE="user.proto";

protoc \
  --go_out=plugins=grpc:"${DIR}${POST_OUT_DIR}" \
  -I".${POST_OUT_DIR}" "${POST_PROTO_FILE}"

protoc \
  --go_out=plugins=grpc:"${DIR}${TAG_OUT_DIR}" \
  -I".${TAG_OUT_DIR}" "${TAG_PROTO_FILE}"

protoc \
  --go_out=plugins=grpc:"${DIR}${USER_OUT_DIR}" \
  -I".${USER_OUT_DIR}" "${USER_PROTO_FILE}"
