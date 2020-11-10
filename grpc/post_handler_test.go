package grpc

import (
	"context"
	"log"
	"testing"

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
		Title:        "Title",
		Content:      "Content",
		MaxNum:       2,
		Gender:       1,
		CreateUserId: 111111,
	})

	createPosts = append(createPosts, &postservice.Post{
		Title:        "Title",
		Content:      "Content",
		MaxNum:       2,
		Gender:       1,
		CreateUserId: 222222,
	})

	createPosts = append(createPosts, &postservice.Post{
		Title:        "Title",
		Content:      "Content",
		MaxNum:       2,
		Gender:       1,
		CreateUserId: 333333,
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
		Title:        "Title",
		Content:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		MaxNum:       2,
		Gender:       1,
		CreateUserId: 555555,
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
	assert.Equal(t, StatusPostContentStringCount, d)
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
		Title:        "Title",
		Content:      "",
		MaxNum:       2,
		Gender:       1,
		CreateUserId: 666666,
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
	assert.Equal(t, StatusPostContentStringCount, d)

}

// TestCreatePostTitleMax Titleが文字数超過の異常系
func TestCreatePostTitleMax(t *testing.T) {
	var createPost = &postservice.Post{
		Title:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Content:      "Content",
		MaxNum:       2,
		Gender:       1,
		CreateUserId: 555555,
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

	assert.Equal(t, "Title", f)
	assert.Equal(t, StatusPostTitleStringCount, d)
}

// TestCreatePostTitleNull Titleが空白の異常系
func TestCreatePostTitleNull(t *testing.T) {
	var createPost = &postservice.Post{
		Title:        "",
		Content:      "Content",
		MaxNum:       2,
		Gender:       1,
		CreateUserId: 555555,
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

	assert.Equal(t, "Title", f)
	assert.Equal(t, StatusPostTitleStringCount, d)
}
