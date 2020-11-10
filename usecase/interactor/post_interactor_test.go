package interactor

import (
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

var DemoPostContentNull = model.Post{
	Title:        testTitle,
	Content:      "",
	MaxNum:       two,
	Gender:       two,
	CreateUserID: testPostUserID,
}

// TestCreate ユーザー作成の正常系
func TestCreate(t *testing.T) {
	// initTable()
	var i PostInteractor
	post := &DemoPost
	createdUser, err := i.Create(post)

	assert.Equal(t, nil, err)
	assert.Equal(t, post.CreateUserID, createdUser.CreateUserID)
	assert.Equal(t, post.Content, createdUser.Content)
	assert.NotEqual(t, 0, createdUser.ID)
}

func TestCreateContentNull(t *testing.T) {
	var i PostInteractor
	post := &DemoPostContentNull
	_, err := i.Create(post)

	assert.NotEqual(t, nil, err)
}

func TestCreateContentTooLong(t *testing.T) {
	var i PostInteractor
	post := &model.Post{CreateUserID: testPostUserID, Title: testTitle, Content: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	_, err := i.Create(post)

	assert.NotEqual(t, nil, err)
}

func TestCreateTitleNull(t *testing.T) {
	var i PostInteractor
	post := &model.Post{
		Title:        "",
		Content:      testContent,
		MaxNum:       two,
		Gender:       two,
		CreateUserID: testUserID,
	}
	_, err := i.Create(post)

	assert.NotEqual(t, nil, err)
}

func TestCreateTitleTooLong(t *testing.T) {
	var i PostInteractor
	post := &model.Post{
		Title:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Content:      testContent,
		MaxNum:       two,
		Gender:       two,
		CreateUserID: testUserID,
	}
	_, err := i.Create(post)

	assert.NotEqual(t, nil, err)
}

func TestDelete(t *testing.T) {
	var i PostInteractor
	post := &model.Post{
		Title:        testTitle,
		Content:      testContent,
		MaxNum:       two,
		Gender:       two,
		CreateUserID: testPostUserID,
	}

	cretedPost, err := i.Create(post)

	assert.Equal(t, nil, err)
	err = i.Delete(cretedPost)
	assert.Equal(t, nil, err)

	deletedPost, err := i.Read(cretedPost.ID)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, zero, deletedPost.ID)
	assert.Equal(t, zero, deletedPost.CreateUserID)
	assert.Equal(t, "", deletedPost.Title)
	assert.Equal(t, "", deletedPost.Content)

}

func TestUpdate(t *testing.T) {
	var i PostInteractor
	post := &model.Post{
		Title:        testTitle,
		Content:      testContent,
		MaxNum:       two,
		Gender:       two,
		CreateUserID: testUserID,
	}
	createdPost, err := i.Create(post)

	assert.Equal(t, nil, err)
	createdAt := createdPost.CreatedAt
	updatePost := createdPost
	updatePost.Content = "Content updated"
	updatePost.UpdateUserID = 12345

	time.Sleep(time.Second * 10)
	updatedPost, err := i.Update(updatePost)

	assert.Equal(t, nil, err)
	readPost, err := i.Read(updatePost.ID)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "content", updatedPost.Content)
	assert.Equal(t, createdPost.ID, updatedPost.ID)
	assert.Equal(t, createdPost.CreateUserID, updatedPost.CreateUserID)
	assert.Equal(t, createdAt, updatedPost.CreatedAt)
	assert.NotEqual(t, readPost.UpdatedAt, updatedPost.UpdatedAt)
	assert.NotEqual(t, testUserID, updatedPost.UpdateUserID)

}

func initTable() {
	DB := db.GetDB()
	DB.Delete(&model.Post{})
}
