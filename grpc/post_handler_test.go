package grpc

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/yzmw1213/PostService/grpc/postservice"
	"github.com/yzmw1213/PostService/grpc/tagservice"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func init() {
	lis = bufconn.Listen(bufSize)
	s := makeServer()
	// 投稿サービス登録
	postservice.RegisterPostServiceServer(s, &server{})
	// タグサービス登録
	tagservice.RegisterTagServiceServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func TestCreatePost(t *testing.T) {
	var createPosts []*postservice.Post
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := postservice.NewPostServiceClient(conn)

	createPosts = append(createPosts, &postservice.Post{
		CreateUserId: 111111,
		Title:        "Title",
		Content:      "Content",
	})

	createPosts = append(createPosts, &postservice.Post{
		CreateUserId: 222222,
		Title:        "Title",
		Content:      "Content",
	})

	createPosts = append(createPosts, &postservice.Post{
		CreateUserId: 333333,
		Title:        "Title",
		Content:      "Content",
	})

	for _, post := range createPosts {
		req := &postservice.CreatePostRequest{
			Post: post,
		}

		_, err = client.CreatePost(ctx, req)
		assert.Equal(t, nil, err)
	}
}

// TestCreatepostContentMax Contentが文字数超過の異常系
func TestCreatepostContentMax(t *testing.T) {
	var createPost = &postservice.Post{
		CreateUserId: 555555,
		Title:        "Title",
		Content:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := postservice.NewPostServiceClient(conn)

	req := &postservice.CreatePostRequest{
		Post: createPost,
	}
	_, err = client.CreatePost(context.Background(), req)

	assert.NotEqual(t, nil, err)

	f, d := getErrorDetail(err)

	assert.Equal(t, "Content", f)
	assert.Equal(t, messageContentMax, d)
}

func getErrorDetail(err error) (string, string) {
	var field string
	var description string
	st, _ := status.FromError(err)
	for _, detail := range st.Details() {
		switch dType := detail.(type) {
		case *errdetails.BadRequest:
			for _, violation := range dType.GetFieldViolations() {
				field = violation.GetField()
				description = violation.GetDescription()
			}
		}
	}

	return field, description
}

// // TestCreatePostContentNull TitleがNullの異常系
func TestCreatePostContentNull(t *testing.T) {
	var createPost = &postservice.Post{
		CreateUserId: 666666,
		Title:        "Title",
		Content:      "",
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := postservice.NewPostServiceClient(conn)

	req := &postservice.CreatePostRequest{
		Post: createPost,
	}
	_, err = client.CreatePost(context.Background(), req)
	assert.NotEqual(t, nil, err)

	f, d := getErrorDetail(err)

	assert.Equal(t, "Content", f)
	assert.Equal(t, messageContentMin, d)

}

func TestDeletePost(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := postservice.NewPostServiceClient(conn)

	createPost := &postservice.Post{
		CreateUserId: 444444,
		Title:        "Title",
		Content:      "Content of the fourth post",
	}

	createReq := &postservice.CreatePostRequest{
		Post: createPost,
	}

	createRes, err := client.CreatePost(ctx, createReq)

	deletePostID := createRes.GetPost().GetId()

	deleteReq := &postservice.DeletePostRequest{
		Id: deletePostID,
	}

	_, err = client.DeletePost(ctx, deleteReq)
	assert.Equal(t, nil, err)

	readReq := &postservice.ReadPostRequest{
		Id: deletePostID,
	}

	readRes, err := client.ReadPost(ctx, readReq)
	assert.NotEqual(t, nil, err)

	_, d := getErrorDetail(err)

	assert.Equal(t, zero, readRes.GetPost().GetId())
	assert.Equal(t, "", readRes.GetPost().GetContent())
	assert.Equal(t, zero, readRes.GetPost().GetCreateUserId())
	assert.Equal(t, "record not found", d)

}

func TestUpdatePost(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	createPost := &postservice.Post{
		CreateUserId: 555555,
		Title:        "Title",
		Content:      "Content",
	}

	client := postservice.NewPostServiceClient(conn)

	createReq := createPostRequest(createPost)

	res, err := client.CreatePost(ctx, createReq)

	updatePost := res.GetPost()
	updatePost.Content = "Content updated"
	time.Sleep(time.Second * 10)

	req := updatePostRequest(updatePost)
	updateRes, err := client.UpdatePost(ctx, req)

	assert.Equal(t, nil, err)

	updatePost = updateRes.GetPost()

	readReq := &postservice.ReadPostRequest{
		Id: updatePost.Id,
	}

	readRes, err := client.ReadPost(context.Background(), readReq)
	assert.Equal(t, nil, err)
	assert.Equal(t, readRes.Post.Id, updatePost.Id)
	assert.Equal(t, updatePost.Content, readRes.Post.Content)

}
