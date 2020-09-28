package interactor

import (
	"log"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/yzmw1213/PostService/domain/model"
)

var DemoValidTag = model.Tag{
	TagName:      testTagName,
	CreateUserID: testUserID,
	Status:       ValidTagStatus,
}

var DemoInvalidTag = model.Tag{
	TagName:      testTagName,
	CreateUserID: testUserID,
	Status:       InValidTagStatus,
}

// TestCreateTag 有効タグ作成の正常系
func TestCreateTag(t *testing.T) {
	var i TagInteractor
	tag := &DemoValidTag
	createdTag, err := i.Create(tag)

	assert.Equal(t, nil, err)
	assert.Equal(t, tag.CreateUserID, createdTag.CreateUserID)
	assert.Equal(t, tag.TagName, createdTag.TagName)
	assert.NotEqual(t, "", createdTag.ID)
}

// TestCreateTag 無効タグ作成の正常系
func TestCreateInvalidTag(t *testing.T) {
	var i TagInteractor
	tag := &DemoInvalidTag
	createdTag, err := i.Create(tag)

	assert.Equal(t, nil, err)
	assert.Equal(t, tag.CreateUserID, createdTag.CreateUserID)
	assert.Equal(t, tag.TagName, createdTag.TagName)
	assert.NotEqual(t, "", createdTag.ID)
}

// TestSearchValidTag 有効タグ検索
func TestSearchValidTag(t *testing.T) {
	var i TagInteractor

	searchdTags, err := i.GetValidTag()

	assert.Equal(t, nil, err)

	for _, tag := range searchdTags {
		log.Println("status")
		log.Println(tag.Status)
		assert.Equal(t, ValidTagStatus, tag.Status)
	}
}

// TestDeleteTag タグ削除
func TestDeleteTag(t *testing.T) {
	var i TagInteractor
	var tagName string = "delete_tag"
	tag := &model.Tag{
		TagName:      tagName,
		CreateUserID: testUserID,
		Status:       1,
	}
	createdTag, err := i.Create(tag)
	assert.Equal(t, nil, err)

	err = i.DeleteByID(createdTag.ID)
	assert.Equal(t, nil, err)

	searchTag, err := i.GetTagByTagName(tagName)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, zero, searchTag.ID)
	assert.Equal(t, "", searchTag.TagName)
}

// TestUpdateTagStatus タグステータス更新
func TestUpdateTagStatus(t *testing.T) {
	var i TagInteractor
	var tagName string = "update_tag"
	tag := &model.Tag{
		TagName:      tagName,
		CreateUserID: testUserID,
		Status:       InValidTagStatus,
	}
	createdTag, err := i.Create(tag)
	assert.Equal(t, nil, err)

	log.Println(createdTag.Status)
	inputTag := createdTag
	inputTag.Status = ValidTagStatus

	_, err = i.Update(inputTag)

	assert.Equal(t, nil, err)
	searchTag, err := i.GetTagByTagName(tagName)
	assert.Equal(t, nil, err)
	assert.Equal(t, tagName, searchTag.TagName)
	assert.Equal(t, ValidTagStatus, searchTag.Status)
}

// TestListTag タグ全件取得
func TestListTag(t *testing.T) {
	var i TagInteractor
	tags, err = i.List()
	log.Println(len(tags))
	assert.Equal(t, nil, err)
}
