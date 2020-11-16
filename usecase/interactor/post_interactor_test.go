package interactor

import (
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

var DemoPostTag = []model.PostTag{
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
	post := &DemoPost
	//
	joinPost := &model.JoinPost{
		Post: &DemoPost,
		Tags: DemoPostTag,
	}
	createdPost, err := i.Create(joinPost)

	assert.Equal(t, nil, err)
	assert.Equal(t, post.CreateUserID, createdPost.Post.CreateUserID)
	assert.Equal(t, post.Content, createdPost.Post.Content)
	assert.NotEqual(t, 0, createdPost.Post.ID)

	// Insertされた各PostTagのPostIDが投稿のIDと等しい事を確認
	// postTags := createdPost.Tags
	// for _, postTag := range postTags {
	// 	assert.Equal(t, createdPost.Post.ID, postTag.PostID)
	// }
}

func TestCreateContentNull(t *testing.T) {
	var i PostInteractor
	joinPost := &model.JoinPost{
		Post: &DemoPostContentNull,
		Tags: DemoPostTag,
	}
	_, err := i.Create(joinPost)
	// PostTagがInsertされていない事をテスト

	assert.NotEqual(t, nil, err)
}

func TestCreateContentTooLong(t *testing.T) {
	var i PostInteractor
	joinPost := &model.JoinPost{
		Post: &model.Post{CreateUserID: testPostUserID, Title: testTitle, Content: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		Tags: DemoPostTag,
	}

	_, err := i.Create(joinPost)

	assert.NotEqual(t, nil, err)
}

func TestCreateTitleNull(t *testing.T) {
	var i PostInteractor
	joinPost := &model.JoinPost{
		Post: &model.Post{
			Title:        "",
			Content:      testContent,
			MaxNum:       two,
			Gender:       two,
			CreateUserID: testUserID,
		},
		Tags: DemoPostTag,
	}

	_, err := i.Create(joinPost)

	assert.NotEqual(t, nil, err)
}

func TestCreateTitleTooLong(t *testing.T) {
	var i PostInteractor
	joinPost := &model.JoinPost{
		Post: &model.Post{
			Title:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Content:      testContent,
			MaxNum:       two,
			Gender:       two,
			CreateUserID: testUserID,
		},
		Tags: DemoPostTag,
	}
	_, err := i.Create(joinPost)

	assert.NotEqual(t, nil, err)
}

func TestDelete(t *testing.T) {
	var i PostInteractor
	joinPost := &model.JoinPost{
		Post: &model.Post{
			Title:        testTitle,
			Content:      testContent,
			MaxNum:       two,
			Gender:       two,
			CreateUserID: testPostUserID,
		},
	}
	cretedPost, err := i.Create(joinPost)

	assert.Equal(t, nil, err)
	err = i.Delete(cretedPost.Post)
	assert.Equal(t, nil, err)

	deletedPost, err := i.Read(cretedPost.Post.ID)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, zero, deletedPost.ID)
	assert.Equal(t, zero, deletedPost.CreateUserID)
	assert.Equal(t, "", deletedPost.Title)
	assert.Equal(t, "", deletedPost.Content)

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

func initTable() {
	DB := db.GetDB()
	DB.Delete(&model.Post{})
	DB.Delete(&model.Tag{})
	DB.Delete(&model.PostTag{})
}
