package grpc

import (
	"context"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/yzmw1213/PostService/grpc/post_grpc"
	"google.golang.org/grpc"
)

func TestCreate(t *testing.T) {
	var createTag *post_grpc.Tag
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := post_grpc.NewTagServiceClient(conn)

	createTag = &post_grpc.Tag{
		TagName:      "tagName",
		Status:       one,
		CreateUserId: "demoUser1",
		UpdateUserId: "",
	}

	createReq := &post_grpc.CreateTagRequest{
		Tag: createTag,
	}

	res, err := client.CreateTag(ctx, createReq)

	assert.Equal(t, nil, err)
	assert.Equal(t, StatusCreateTagSuccess, res.GetStatus().GetCode())
}

func TestCreateTagNameNull(t *testing.T) {
	var createTag *post_grpc.Tag
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := post_grpc.NewTagServiceClient(conn)

	createTag = &post_grpc.Tag{
		TagName:      "",
		Status:       one,
		CreateUserId: "demoUser1",
		UpdateUserId: "",
	}

	createReq := &post_grpc.CreateTagRequest{
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
	client := post_grpc.NewTagServiceClient(conn)

	req := &post_grpc.ListTagRequest{}

	_, err = client.ListTag(ctx, req)

	assert.Equal(t, nil, err)
}
