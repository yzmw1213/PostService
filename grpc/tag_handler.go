package grpc

import (
	"context"
	"log"

	"github.com/yzmw1213/PostService/domain/model"
	"github.com/yzmw1213/PostService/grpc/tagservice"
)

const (
	// StatusCreateTagSuccess タグ作成成功ステータス
	StatusCreateTagSuccess string = "TAG_CREATE_SUCCESS"
	// StatusUpdateTagSuccess タグ更新作成成功ステータス
	StatusUpdateTagSuccess string = "TAG_UPDATE_SUCCESS"
	// StatusDeleteTagSuccess タグ削除成功ステータス
	StatusDeleteTagSuccess string = "TAG_DELETE_SUCCESS"
	// StatusTagNotExists 指定したタグの登録がない時のエラーステータス
	StatusTagNotExists string = "TAG_NOT_EXISTS_ERROR"
	// StatustagNameAlreadyUsed 既に使われているTagName登録時のエラーステータス
	StatustagNameAlreadyUsed string = "TAG_NAME_ALREADY_USED_ERROR"
	// StatusTagNameStringCount タグ名文字数が無効のエラーステータス
	StatusTagNameStringCount string = "TAG_NAME_COUNT_ERROR"
)

func (s server) CreateTag(ctx context.Context, req *tagservice.CreateTagRequest) (*tagservice.CreateTagResponse, error) {
	postData := req.GetTag()
	tag := makeTagModel(postData)

	// 既に同一のtagnameによる登録がないかチェック
	if s.tagExistsByTagName(tag.TagName) == true {
		return s.makeCreateTagResponse(StatustagNameAlreadyUsed), nil
	}

	tag, err := s.TagUsecase.Create(tag)
	if err != nil {
		return nil, err
	}

	return s.makeCreateTagResponse(StatusCreateTagSuccess), nil
}

func (s server) DeleteTag(ctx context.Context, req *tagservice.DeleteTagRequest) (*tagservice.DeleteTagResponse, error) {
	id := req.GetTagId()

	// 既にタグが削除されていないかチェック
	if s.tagExistsByTagID(id) != true {
		log.Println("tag not exists")
		return s.makeDeleteTagResponse(StatusTagNotExists), nil
	}

	err = s.TagUsecase.DeleteByID(id)
	if err != nil {
		return nil, err
	}

	return s.makeDeleteTagResponse(StatusDeleteTagSuccess), nil
}

func (s server) UpdateTag(ctx context.Context, req *tagservice.UpdateTagRequest) (*tagservice.UpdateTagResponse, error) {
	postData := req.GetTag()

	tag := makeTagModel(postData)

	if _, err := s.TagUsecase.Update(tag); err != nil {
		return nil, err
	}

	return s.makeUpdateTagResponse(StatusUpdateTagSuccess), nil
}

// ListTag 全てのタグを取得して返す
func (s server) ListTag(ctx context.Context, req *tagservice.ListTagRequest) (*tagservice.ListTagResponse, error) {
	rows, err := s.TagUsecase.List()
	if err != nil {
		return nil, err
	}
	var tags []*tagservice.Tag
	for _, tag := range rows {
		tag := makeGrpcTag(&tag)
		tags = append(tags, tag)
	}
	res := &tagservice.ListTagResponse{
		Tag: tags,
	}
	return res, nil
}

// ListValidTag 公開ステータスが公開のタグを取得して返す
func (s server) ListValidTag(ctx context.Context, req *tagservice.ListValidTagRequest) (*tagservice.ListValidTagResponse, error) {
	rows, err := s.TagUsecase.ListValidTag()
	if err != nil {
		return nil, err
	}
	var tags []*tagservice.Tag
	for _, tag := range rows {
		tag := makeGrpcTag(&tag)
		tags = append(tags, tag)
	}
	res := &tagservice.ListValidTagResponse{
		Tag: tags,
	}
	return res, nil
}

func makeTagModel(gTag *tagservice.Tag) *model.Tag {
	tag := &model.Tag{
		ID:           gTag.GetTagId(),
		TagName:      gTag.GetTagName(),
		Status:       gTag.GetStatus(),
		CreateUserID: gTag.GetCreateUserId(),
		UpdateUserID: gTag.GetUpdateUserId(),
	}
	return tag
}

func makeGrpcTag(tag *model.Tag) *tagservice.Tag {
	gTag := &tagservice.Tag{
		TagId:        tag.ID,
		TagName:      tag.TagName,
		Status:       tag.Status,
		CreateUserId: tag.CreateUserID,
		UpdateUserId: tag.UpdateUserID,
	}
	return gTag
}

// makeCreateTagResponse CreateTagメソッドのresponseを生成し返す
func (s server) makeCreateTagResponse(statusCode string) *tagservice.CreateTagResponse {
	res := &tagservice.CreateTagResponse{}
	if statusCode != "" {
		responseStatus := &tagservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeDeleteTagResponse CreateTagメソッドのresponseを生成し返す
func (s server) makeDeleteTagResponse(statusCode string) *tagservice.DeleteTagResponse {
	res := &tagservice.DeleteTagResponse{}
	if statusCode != "" {
		responseStatus := &tagservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeUpdateTagResponse UpdateTagメソッドのresponseを生成し返す
func (s server) makeUpdateTagResponse(statusCode string) *tagservice.UpdateTagResponse {
	res := &tagservice.UpdateTagResponse{}
	if statusCode != "" {
		responseStatus := &tagservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// tagExistsByTagName 同名のタグが登録済みかの判定
func (s server) tagExistsByTagName(tagName string) bool {
	if tagName == "" {
		return false
	}
	tag, _ := s.TagUsecase.GetTagByTagName(tagName)
	if tag.ID == 0 {
		return false
	}
	return true
}

// tagExistsByTagID　IDが一致するタグの登録があるかの判定
func (s server) tagExistsByTagID(tagID uint32) bool {
	tag, _ := s.TagUsecase.GetTagByTagID(tagID)
	if tag.ID == 0 {
		return false
	}
	return true
}
