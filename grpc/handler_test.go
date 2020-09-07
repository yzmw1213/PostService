package grpc

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/yzmw1213/PostService/grpc/post_grpc"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var err error

func init() {
	lis = bufconn.Listen(bufSize)
	s := makeServer()
	post_grpc.RegisterPostServiceServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreatePost(t *testing.T) {
	var createPosts []*post_grpc.Post
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := post_grpc.NewPostServiceClient(conn)

	createPosts = append(createPosts, &post_grpc.Post{
		UserId:  111111,
		Content: "Content",
	})

	createPosts = append(createPosts, &post_grpc.Post{
		UserId:  222222,
		Content: "Content",
	})

	createPosts = append(createPosts, &post_grpc.Post{
		UserId:  333333,
		Content: "Content",
	})

	for _, post := range createPosts {
		req := &post_grpc.CreatePostRequest{
			Post: post,
		}

		_, err = client.CreatePost(ctx, req)
		assert.Equal(t, nil, err)
	}
}

// TestCreatepostContentMax Contentが文字数超過の異常系
func TestCreatepostContentMax(t *testing.T) {
	var createPost = &post_grpc.Post{
		UserId:  555555,
		Content: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := post_grpc.NewPostServiceClient(conn)

	req := &post_grpc.CreatePostRequest{
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
	log.Println("getErrorDetail")
	log.Println(field)
	log.Println(description)

	return field, description
}

// // TestCreatePostContentNull TitleがNullの異常系
func TestCreatePostContentNull(t *testing.T) {
	var createPost = &post_grpc.Post{
		UserId:  666666,
		Content: "",
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := post_grpc.NewPostServiceClient(conn)

	req := &post_grpc.CreatePostRequest{
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

	client := post_grpc.NewPostServiceClient(conn)

	createPost := &post_grpc.Post{
		UserId:  444444,
		Content: "Content of the fourth post",
	}

	createReq := &post_grpc.CreatePostRequest{
		Post: createPost,
	}

	createRes, err := client.CreatePost(ctx, createReq)

	deletePostId := createRes.GetPost().GetId()

	deleteReq := &post_grpc.DeletePostRequest{
		Id: deletePostId,
	}

	_, err = client.DeletePost(ctx, deleteReq)
	assert.Equal(t, nil, err)

	readReq := &post_grpc.ReadPostRequest{
		Id: deletePostId,
	}

	readRes, err := client.ReadPost(ctx, readReq)
	assert.NotEqual(t, nil, err)

	_, d := getErrorDetail(err)

	assert.Equal(t, int32(0), readRes.GetPost().GetId())
	assert.Equal(t, "", readRes.GetPost().GetContent())
	assert.Equal(t, int32(0), readRes.GetPost().GetUserId())
	assert.Equal(t, "record not found", d)

}

func TestUpdatePost(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	createPost := &post_grpc.Post{
		UserId:  555555,
		Content: "Content",
	}

	client := post_grpc.NewPostServiceClient(conn)

	createReq := createPostRequest(createPost)

	res, err := client.CreatePost(ctx, createReq)

	updatePost := res.GetPost()
	updatePost.Content = "Content updated"
	time.Sleep(time.Second * 10)

	req := updatePostRequest(updatePost)
	updateRes, err := client.UpdatePost(ctx, req)

	assert.Equal(t, nil, err)

	updatePost = updateRes.GetPost()

	readReq := &post_grpc.ReadPostRequest{
		Id: updatePost.Id,
	}

	readRes, err := client.ReadPost(context.Background(), readReq)
	assert.Equal(t, nil, err)
	assert.Equal(t, readRes.Post.Id, updatePost.Id)
	assert.Equal(t, updatePost.Content, readRes.Post.Content)

}
