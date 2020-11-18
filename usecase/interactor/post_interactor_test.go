package interactor

import (
	"log"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/yzmw1213/PostService/db"
	"github.com/yzmw1213/PostService/domain/model"
)

var DemoPost = model.Post{
	Title:        testTitle,
	Content:      testContent,
	MaxNum:       two,
	Gender:       two,
	CreateUserID: testPostUserID,
}

var DemoPostTag = model.PostTag{}

var DemoPostTags = []model.PostTag{
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

var DemoPostContentNull = model.Post{
	Title:        testTitle,
	Content:      "",
	MaxNum:       two,
	Gender:       two,
	CreateUserID: testPostUserID,
}

var DemoJoinPost = model.JoinPost{}

// TestCreate 投稿作成の正常系
func TestCreate(t *testing.T) {
	initTable()
	var i PostInteractor
	post := makePost(testTitle, testContent, two, two)
	postTags := makePostTags()

	//
	joinPost := createDemoJoinPost(post, postTags)

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
	post := makePost(testTitle, "", two, two)
	postTags := makePostTags()
	joinPost := createDemoJoinPost(post, postTags)

	beforePostTagCount := countPostTag()
	_, err := i.Create(&joinPost)
	// PostTagがInsertされていない事をテスト
	afterPostTagCount := countPostTag()
	assert.Equal(t, beforePostTagCount, afterPostTagCount)

	assert.NotEqual(t, nil, err)
}

func TestCreateContentTooLong(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", two, two)
	postTags := makePostTags()
	joinPost := createDemoJoinPost(post, postTags)

	beforePostTagCount := countPostTag()

	_, err := i.Create(&joinPost)
	afterPostTagCount := countPostTag()

	assert.Equal(t, beforePostTagCount, afterPostTagCount)
	assert.NotEqual(t, nil, err)
}

func TestCreateTitleNull(t *testing.T) {
	var i PostInteractor
	post := makePost("", testContent, two, two)
	postTags := makePostTags()
	joinPost := createDemoJoinPost(post, postTags)

	beforePostTagCount := countPostTag()

	_, err := i.Create(&joinPost)
	afterPostTagCount := countPostTag()

	assert.NotEqual(t, nil, err)
	assert.Equal(t, beforePostTagCount, afterPostTagCount)
}

func TestCreateTitleTooLong(t *testing.T) {
	var i PostInteractor
	post := makePost("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", testContent, two, two)
	postTags := makePostTags()
	joinPost := createDemoJoinPost(post, postTags)
	beforePostTagCount := countPostTag()

	_, err := i.Create(&joinPost)
	afterPostTagCount := countPostTag()

	assert.NotEqual(t, nil, err)
	assert.Equal(t, beforePostTagCount, afterPostTagCount)
}

func TestDelete(t *testing.T) {
	var i PostInteractor
	log.Println("DemoPost", DemoPost)
	post := makePost(testTitle, testContent, two, two)
	postTags := makePostTags()
	joinPost := createDemoJoinPost(post, postTags)

	cretedJoinPost, err := i.Create(&joinPost)

	assert.Equal(t, nil, err)
	postID := cretedJoinPost.Post.ID
	beforePostTagCount := countPostTag()

	err = i.DeleteByID(postID)
	assert.Equal(t, nil, err)

	log.Println("created PostID", postID)

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
}

// func TestUpdate(t *testing.T) {
// 	var i PostInteractor
// 	post := &model.Post{
// 		Title:        testTitle,
// 		Content:      testContent,
// 		MaxNum:       two,
// 		Gender:       two,
// 		CreateUserID: testUserID,
// 	}
// 	createdPost, err := i.Create(post)

// 	assert.Equal(t, nil, err)
// 	createdAt := createdPost.CreatedAt
// 	updatePost := createdPost
// 	updatePost.Content = "Content updated"
// 	updatePost.UpdateUserID = 12345

// 	time.Sleep(time.Second * 10)
// 	updatedPost, err := i.Update(updatePost)

// 	assert.Equal(t, nil, err)
// 	readPost, err := i.Read(updatePost.ID)
// 	assert.Equal(t, nil, err)
// 	assert.NotEqual(t, "content", updatedPost.Content)
// 	assert.Equal(t, createdPost.ID, updatedPost.ID)
// 	assert.Equal(t, createdPost.CreateUserID, updatedPost.CreateUserID)
// 	assert.Equal(t, createdAt, updatedPost.CreatedAt)
// 	assert.NotEqual(t, readPost.UpdatedAt, updatedPost.UpdatedAt)
// 	assert.NotEqual(t, testUserID, updatedPost.UpdateUserID)

// }

// func TestUpdate(t *testing.T) {
// 	var i PostInteractor
// 	post := makePost(testTitle, testContent, two, two)
// 	postTags := makePostTags()
// 	joinPost := createDemoJoinPost(post, postTags)
// 	createdJoinPost, err := i.Create(&joinPost)

// 	assert.Equal(t, nil, err)
// 	createdAt := createdJoinPost.Post.CreatedAt
// 	updateJoinPost := createdJoinPost
// 	updateJoinPost.Post.Content = "Content updated"
// 	updateJoinPost.Post.UpdateUserID = 12345

// 	time.Sleep(time.Second * 3)
// 	updatedPost, err := i.Update(&updateJoinPost)

// 	assert.Equal(t, nil, err)
// 	readPost, err := i.Read(updatePost.ID)
// 	assert.Equal(t, nil, err)
// 	assert.NotEqual(t, "content", updatedPost.Content)
// 	assert.Equal(t, createdPost.ID, updatedPost.ID)
// 	assert.Equal(t, createdPost.CreateUserID, updatedPost.CreateUserID)
// 	assert.Equal(t, createdAt, updatedPost.CreatedAt)
// 	assert.NotEqual(t, readPost.UpdatedAt, updatedPost.UpdatedAt)
// 	assert.NotEqual(t, testUserID, updatedPost.UpdateUserID)

// }

func createDemoJoinPost(post model.Post, postTags []model.PostTag) model.JoinPost {
	joinPost := model.JoinPost{
		Post:     &post,
		PostTags: postTags,
	}
	return joinPost
}

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

func makePost(title string, content string, maxNum uint32, gender uint32) model.Post {
	return model.Post{
		ID:           zero,
		Title:        title,
		Content:      content,
		MaxNum:       maxNum,
		Gender:       gender,
		CreateUserID: testPostUserID,
		UpdateUserID: zero,
	}
}

func initTable() {
	DB := db.GetDB()
	DB.Delete(&model.Post{})
	DB.Delete(&model.Tag{})
	DB.Delete(&model.PostTag{})
}