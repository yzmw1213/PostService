package grpc

import (
	"context"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/yzmw1213/PostService/grpc/tagservice"
	"google.golang.org/grpc"
)

func TestCreate(t *testing.T) {
	var createTag *tagservice.Tag
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := tagservice.NewTagServiceClient(conn)

	createTag = &tagservice.Tag{
		TagName:      "tagName",
		Status:       one,
		CreateUserId: 1,
		UpdateUserId: 0,
	}

	createReq := &tagservice.CreateTagRequest{
		Tag: createTag,
	}

	res, err := client.CreateTag(ctx, createReq)

	assert.Equal(t, nil, err)
	assert.Equal(t, StatusCreateTagSuccess, res.GetStatus().GetCode())
}

func TestCreateTagNameNull(t *testing.T) {
	var createTag *tagservice.Tag
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := tagservice.NewTagServiceClient(conn)

	createTag = &tagservice.Tag{
		TagName:      "",
		Status:       one,
		CreateUserId: 1,
		UpdateUserId: 0,
	}

	createReq := &tagservice.CreateTagRequest{
		Tag: createTag,
	}

	_, err = client.CreateTag(ctx, createReq)

	f, d := getErrorDetail(err)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "TagName", f)
	assert.Equal(t, StatusTagNameStringCount, d)
}

func TestList(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := tagservice.NewTagServiceClient(conn)

	req := &tagservice.ListTagRequest{}

	_, err = client.ListTag(ctx, req)

	assert.Equal(t, nil, err)
}
