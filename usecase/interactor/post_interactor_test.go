package interactor

import (
	"log"
	"testing"
	"time"

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
	joinPost := makeJoinPost(post, postTags)

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
	joinPost := makeJoinPost(post, postTags)

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
	joinPost := makeJoinPost(post, postTags)

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
	joinPost := makeJoinPost(post, postTags)

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
	joinPost := makeJoinPost(post, postTags)
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
	joinPost := makeJoinPost(post, postTags)

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

func TestUpdatePost(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, testContent, two, two)
	joinPost := makeJoinPost(post, nil)
	createdJoinPost, err := i.Create(&joinPost)

	assert.Equal(t, nil, err)
	updateJoinPost := createdJoinPost
	updateJoinPost.Post.Content = "Content updated"
	updateJoinPost.Post.UpdateUserID = 12345

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
	assert.NotEqual(t, testPostUserID, updatedPost.UpdateUserID)

}

func TestUpdatePostTag(t *testing.T) {
	var i PostInteractor
	post := makePost(testTitle, testContent, two, two)
	postTags := makePostTags()
	joinPost := makeJoinPost(post, postTags)

	createdPost, err := i.Create(&joinPost)
	assert.Equal(t, nil, err)
	postID := createdPost.Post.ID
	beforePostTagCount := countPostTagByPostID(postID)

	updatePostTags := makePostTags()
	updatePostTags = append(updatePostTags, model.PostTag{PostID: postID, TagID: four})

	joinPost.PostTags = updatePostTags

	_, err = i.Update(&joinPost)

	afterPostTagCount := countPostTagByPostID(postID)

	assert.Equal(t, beforePostTagCount+1, afterPostTagCount)
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
