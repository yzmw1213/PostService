package interactor

import (
	"log"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/yzmw1213/PostService/db"
	"github.com/yzmw1213/PostService/domain/model"
)

var (
	DemoPost            = model.Post{Title: testTitle, Content: testContent, CreateUserID: testPostUserID}
	DemoPostTag         = model.PostTag{}
	DemoPostContentNull = model.Post{Title: testTitle, Content: "", CreateUserID: testPostUserID}
	DemoJoinPost        = model.JoinPost{}
	DemoUser            = model.User{}
	DemoPostLikeUser    = []model.User{}
	user1               uint32
	user2               uint32
	user3               uint32
)

// TestCreate 投稿作成の正常系
func TestCreate(t *testing.T) {
	initTable()
	var i PostInteractor
	post := makePost(testTitle, testContent)
	post.CreateUserID = user1
	postTags := makePostTags()

	//
	joinPost := makeJoinPost(post, DemoUser, postTags, nil, nil)

	// 登録前のpostTag登録数
	beforePostTagCount := countPostTag()
	createdPost, err := i.Create(&joinPost)

	assert.Equal(t, nil, err)
	assert.Equal(t, post.CreateUserID, createdPost.Post.CreateUserID)
	assert.Equal(t, post.Content, createdPost.Post.Content)
	assert.NotEqual(t, 0, createdPost.Post.ID)

	// 登録後のpostTag登録数
	afterPostTagCount := countPostTag()

	// Postと同一PostIDのPostTagが登録されている事を確認
	postID := createdPost.Post.ID
	postTags, err = listPostTagsByID(postID)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(postTags))
	assert.NotEqual(t, beforePostTagCount, afterPostTagCount)
}

func TestCreateContentNull(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, "")
	postTags := makePostTags()
	joinPost := makeJoinPost(post, DemoUser, postTags, nil, nil)

	beforePostTagCount := countPostTag()
	_, err := i.Create(&joinPost)
	// PostTagがInsertされていない事をテスト
	afterPostTagCount := countPostTag()
	assert.Equal(t, beforePostTagCount, afterPostTagCount)

	assert.NotEqual(t, nil, err)
}

func TestCreateContentTooLong(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	postTags := makePostTags()
	joinPost := makeJoinPost(post, DemoUser, postTags, nil, nil)

	beforePostTagCount := countPostTag()

	_, err := i.Create(&joinPost)
	afterPostTagCount := countPostTag()

	assert.Equal(t, beforePostTagCount, afterPostTagCount)
	assert.NotEqual(t, nil, err)
}

func TestCreateTitleNull(t *testing.T) {
	var i PostInteractor
	post := makePost("", testContent)
	postTags := makePostTags()
	joinPost := makeJoinPost(post, DemoUser, postTags, nil, nil)

	beforePostTagCount := countPostTag()

	_, err := i.Create(&joinPost)
	afterPostTagCount := countPostTag()

	assert.NotEqual(t, nil, err)
	assert.Equal(t, beforePostTagCount, afterPostTagCount)
}

func TestCreateTitleTooLong(t *testing.T) {
	var i PostInteractor
	post := makePost("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", testContent)
	postTags := makePostTags()
	joinPost := makeJoinPost(post, DemoUser, postTags, nil, nil)
	beforePostTagCount := countPostTag()

	_, err := i.Create(&joinPost)
	afterPostTagCount := countPostTag()

	assert.NotEqual(t, nil, err)
	assert.Equal(t, beforePostTagCount, afterPostTagCount)
}

func TestDelete(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, testContent)
	post.CreateUserID = user2
	postTags := makePostTags()
	joinPost := makeJoinPost(post, DemoUser, postTags, nil, nil)

	cretedJoinPost, err := i.Create(&joinPost)

	assert.Equal(t, nil, err)
	postID := cretedJoinPost.Post.ID
	beforePostTagCount := countPostTag()

	err = i.DeleteByID(postID)
	assert.Equal(t, nil, err)

	deletedPost, err := i.GetByID(postID)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, zero, deletedPost.ID)
	assert.Equal(t, zero, deletedPost.CreateUserID)
	assert.Equal(t, "", deletedPost.Title)
	assert.Equal(t, "", deletedPost.Content)

	// Postと同一IDのPostTagが全て削除されている事を確認
	postTags, err = listPostTagsByID(postID)
	afterPostTagCount := countPostTag()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(postTags))
	assert.NotEqual(t, beforePostTagCount, afterPostTagCount)

	// Postと同一IDCommentが削除されていない事を確認
}

func TestUpdatePost(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, testContent)
	post.CreateUserID = user2
	joinPost := makeJoinPost(post, DemoUser, nil, nil, nil)
	createdJoinPost, err := i.Create(&joinPost)

	assert.Equal(t, nil, err)
	updateJoinPost := createdJoinPost
	updateJoinPost.Post.Content = "Content updated"
	updateJoinPost.Post.UpdateUserID = user2

	time.Sleep(time.Second * 3)
	updatedJoinPost, err := i.Update(updateJoinPost)
	updatedPost := updatedJoinPost.Post

	assert.Equal(t, nil, err)
	readPost, err := i.GetByID(updatedJoinPost.Post.ID)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "content", updatedPost.Content)
	assert.Equal(t, createdJoinPost.Post.ID, updatedPost.ID)
	assert.Equal(t, createdJoinPost.Post.CreateUserID, updatedPost.CreateUserID)
	assert.Equal(t, createdJoinPost.Post.CreatedAt, updatedPost.CreatedAt)
	assert.NotEqual(t, readPost.UpdatedAt, updatedPost.UpdatedAt)
	assert.Equal(t, user2, updatedPost.UpdateUserID)

}

func TestUpdatePostTag(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, testContent)
	postTags := makePostTags()
	post.CreateUserID = user3
	joinPost := makeJoinPost(post, DemoUser, postTags, nil, nil)

	createdPost, err := i.Create(&joinPost)
	assert.Equal(t, nil, err)
	postID := createdPost.Post.ID
	beforePostTagCount := countPostTagByPostID(postID)

	updatePostTags := makePostTags()
	updatePostTags = append(updatePostTags, model.PostTag{PostID: postID, TagID: four})

	joinPost.PostTags = updatePostTags
	joinPost.Post.UpdateUserID = user3

	updatedJoinPost, err := i.Update(&joinPost)

	afterPostTagCount := countPostTagByPostID(postID)

	assert.Equal(t, beforePostTagCount+1, afterPostTagCount)
	assert.Equal(t, user3, updatedJoinPost.Post.UpdateUserID)
}

func selectUsers() (uint32, uint32, uint32) {
	users := getUserData()
	var num int = 1

	var user1 uint32
	var user2 uint32
	var user3 uint32
	for i := range users {
		if num == 1 {
			user1 = i
			log.Println("user1", user1)
		}
		if num == 2 {
			user2 = i
			log.Println("user2", user2)
		}
		if num == 3 {
			user3 = i
			log.Println("user3", user3)
		}
		num++
	}
	return user1, user2, user3
}

func TestLikePost(t *testing.T) {
	var i PostInteractor

	post := makePost(testTitle, testContent)
	post.CreateUserID = user1
	joinPost := makeJoinPost(post, DemoUser, nil, nil, nil)
	createdPost, err := i.Create(&joinPost)
	assert.Equal(t, nil, err)
	postID := createdPost.Post.ID
	likeUser := &model.PostLikeUser{PostID: postID, UserID: user1}

	_, err = i.Like(likeUser)
	assert.Equal(t, nil, err)

	// likeしているユーザー数をカウントするテスト
	likeCount := countPostLikeUserByPostID(postID)
	assert.Equal(t, 1, likeCount)

	likeUser = &model.PostLikeUser{PostID: postID, UserID: user2}
	_, err = i.Like(likeUser)

	// likeしているユーザー数が増えている事をテスト
	likeCount = countPostLikeUserByPostID(postID)
	assert.Equal(t, 2, likeCount)
}

func TestNotLikePost(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, testContent)
	post.CreateUserID = user2
	joinPost := makeJoinPost(post, DemoUser, nil, nil, nil)
	createdPost, err := i.Create(&joinPost)
	assert.Equal(t, nil, err)
	postID := createdPost.Post.ID
	likeUsers := []model.PostLikeUser{
		{PostID: postID, UserID: user1},
		{PostID: postID, UserID: user2},
		{PostID: postID, UserID: user3},
	}

	for _, user := range likeUsers {
		_, err = i.Like(&user)
	}
	assert.Equal(t, nil, err)

	beforeLikeCount := countPostLikeUserByPostID(postID)

	// お気に入りを1件削除
	_, err = i.NotLike(&model.PostLikeUser{PostID: postID, UserID: user1})
	afterLikeCount := countPostLikeUserByPostID(postID)

	// likeしているユーザー数が1だけ減っている事をテスト
	assert.Equal(t, afterLikeCount, beforeLikeCount-1)
}

func TestCreateComment(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, testContent)
	post.CreateUserID = user3
	joinPost := makeJoinPost(post, DemoUser, nil, nil, nil)
	createdPost, err := i.Create(&joinPost)
	assert.Equal(t, nil, err)

	postID := createdPost.Post.ID
	comment := makeComment(*createdPost.Post, testCommentContent)
	_, err = i.CreateComment(&comment)

	assert.Equal(t, nil, err)

	count := countCommentByPostID(postID)

	assert.Equal(t, 1, count)

}

func TestCreateCommentNull(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, testContent)
	post.CreateUserID = user1
	joinPost := makeJoinPost(post, DemoUser, nil, nil, nil)
	createdPost, err := i.Create(&joinPost)
	assert.Equal(t, nil, err)

	postID := createdPost.Post.ID
	comment := makeComment(*createdPost.Post, "")
	_, err = i.CreateComment(&comment)
	count := countCommentByPostID(postID)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, 0, count)
}

func TestCreateCommentTooLong(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, testContent)
	post.CreateUserID = user1
	joinPost := makeJoinPost(post, DemoUser, nil, nil, nil)
	createdPost, err := i.Create(&joinPost)
	assert.Equal(t, nil, err)

	postID := createdPost.Post.ID
	comment := makeComment(*createdPost.Post, "コメントが入りますコメントが入りますコメントが入りますコメントが入りますコメントが入りますコメントが入りますコメントが入りますコメントが入りますコメントが入りますコメントが入りますコメントが入りますコメントが入りますコメントが入りますコメントが入ります")
	_, err = i.CreateComment(&comment)
	count := countCommentByPostID(postID)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, 0, count)
}

func TestUpdateComment(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, testContent)
	post.CreateUserID = user2
	joinPost := makeJoinPost(post, DemoUser, nil, nil, nil)
	createdPost, err := i.Create(&joinPost)

	assert.Equal(t, nil, err)
	comment := makeComment(*createdPost.Post, testCommentContent)
	createdComment, err := i.CreateComment(&comment)
	updatedAt := createdComment.UpdatedAt

	assert.Equal(t, nil, err)

	updateComment := createdComment
	createdComment.CommentContent = "updated comment content"

	time.Sleep(time.Second * 3)
	updatedComment, err := i.UpdateComment(updateComment)

	readComment, err := getCommentByID(updatedComment.CommentID)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, testCommentContent, readComment.CommentContent)
	assert.Equal(t, createdComment.CreatedAt, updatedComment.CreatedAt)
	assert.NotEqual(t, updatedAt, updatedComment.UpdatedAt)
}

func TestDeleteComment(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, testContent)
	post.CreateUserID = user3
	joinPost := makeJoinPost(post, DemoUser, nil, nil, nil)
	createdPost, err := i.Create(&joinPost)
	assert.Equal(t, nil, err)
	postID := createdPost.Post.ID

	comment := makeComment(*createdPost.Post, testCommentContent)

	createdComment, err := i.CreateComment(&comment)
	assert.Equal(t, nil, err)
	beforeCommentCount := countCommentByPostID(postID)

	assert.Equal(t, 0, beforeCommentCount)
	commentID := createdComment.CommentID

	// コメントを1件削除
	err = i.DeleteComment(commentID)
	assert.Equal(t, nil, err)

	afterCommentCount := countCommentByPostID(postID)
	// likeしているユーザー数が1だけ減っている事をテスト
	assert.Equal(t, 0, afterCommentCount)
}

// func TestGetAllPosts(t *testing.T) {
// 	var i PostInteractor
// 	posts, err := i.List("all", 0)
// 	assert.Equal(t, nil, err)
// 	assert.NotEqual(t, 0, len(posts))
// }

// func TestGetPostsByUserID(t *testing.T) {
// 	var i PostInteractor
// 	posts, err := i.List("create", user1)
// 	assert.Equal(t, nil, err)
// 	assert.NotEqual(t, 0, len(posts))

// 	for _, post := range posts {
// 		assert.Equal(t, user1, post.Post.CreateUserID)
// 	}
// }

// func TestGetPostsByLikeUserID(t *testing.T) {
// 	var i PostInteractor
// 	posts, err := i.List("like", user2)
// 	assert.Equal(t, nil, err)
// 	assert.NotEqual(t, 0, len(posts))
// }

// func TestGetPostsByTagID(t *testing.T) {
// 	var i PostInteractor

// 	posts, err := i.List("tag", one)
// 	assert.Equal(t, nil, err)
// 	assert.NotEqual(t, 0, len(posts))

// 	posts, err = i.List("tag", four)
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, 1, len(posts))
// }

func makePostTags() []model.PostTag {
	return []model.PostTag{
		{
			TagID: one,
		},
		{
			TagID: two,
		},
		{
			TagID: three,
		},
	}
}

func makePost(title string, content string) model.Post {
	return model.Post{
		ID:           zero,
		Title:        title,
		Content:      content,
		CreateUserID: testPostUserID,
		UpdateUserID: zero,
	}
}

func makeComment(post model.Post, comment string) model.Comment {
	return model.Comment{
		PostID:         post.ID,
		CreateUserID:   user1,
		CommentContent: comment,
	}
}

func initTable() {
	DB := db.GetDB()
	// circleCIではリクエストを実行できていないためコメントアウト
	// user1, user2, user3 = selectUsers()
	user1, user2, user3 = one, two, three
	DB.Delete(&model.Post{})
	DB.Delete(&model.Tag{})
	DB.Delete(&model.PostTag{})
	DB.Delete(&model.PostLikeUser{})
}
