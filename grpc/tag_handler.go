package grpc

import (
	"context"
	"log"

	"github.com/yzmw1213/PostService/domain/model"
	"github.com/yzmw1213/PostService/grpc/post_grpc"
)

const (
	// StatusCreateTagSuccess タグ作成成功ステータス
	StatusCreateTagSuccess string = "TAG_CREATE_SUCCESS"
	// StatusUpdateTagSuccess タグ作成成功ステータス
	StatusUpdateTagSuccess string = "TAG_UPDATE_SUCCESS"
	// StatustagNameAlreadyUsed 既に使われているTagName登録時のエラーステータス
	StatustagNameAlreadyUsed string = "TAG_NAME_ALREADY_USED_ERROR"
)

func (s server) CreateTag(ctx context.Context, req *post_grpc.CreateTagRequest) (*post_grpc.CreateTagResponse, error) {
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

func (s server) UpdateTag(ctx context.Context, req *post_grpc.UpdateTagRequest) (*post_grpc.UpdateTagResponse, error) {
	postData := req.GetTag()

	tag := makeTagModel(postData)

	if _, err := s.TagUsecase.Update(tag); err != nil {
		return nil, err
	}

	return s.makeUpdateTagResponse(StatusUpdateTagSuccess), nil
}

func (s server) ListTag(req *post_grpc.ListTagRequest, stream post_grpc.TagService_ListTagServer) error {
	rows, err := s.TagUsecase.List()
	if err != nil {
		return err
	}
	for _, tag := range rows {
		tag := makeGrpcTag(&tag)
		res := &post_grpc.ListTagResponse{
			Tag: tag,
		}
		sendErr := stream.Send(res)
		if sendErr != nil {
			log.Fatalf("Error while sending response to client :%v", sendErr)
			return sendErr
		}
	}

	return nil
}

func makeTagModel(gTag *post_grpc.Tag) *model.Tag {
	tag := &model.Tag{
		ID:           gTag.GetTagId(),
		TagName:      gTag.GetTagName(),
		Status:       gTag.GetStatus(),
		CreateUserID: gTag.GetCreateUserId(),
		UpdateUserID: gTag.GetUpdateUserId(),
	}
	return tag
}

func makeGrpcTag(tag *model.Tag) *post_grpc.Tag {
	gTag := &post_grpc.Tag{
		TagId:        tag.ID,
		TagName:      tag.TagName,
		Status:       tag.Status,
		CreateUserId: tag.CreateUserID,
		UpdateUserId: tag.UpdateUserID,
	}
	return gTag
}

// makeCreateTagResponse CreateTagメソッドのresponseを生成し返す
func (s server) makeCreateTagResponse(statusCode string) *post_grpc.CreateTagResponse {
	res := &post_grpc.CreateTagResponse{}
	if statusCode != "" {
		responseStatus := &post_grpc.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeUpdateTagResponse UpdateTagメソッドのresponseを生成し返す
func (s server) makeUpdateTagResponse(statusCode string) *post_grpc.UpdateTagResponse {
	res := &post_grpc.UpdateTagResponse{}
	if statusCode != "" {
		responseStatus := &post_grpc.ResponseStatus{
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
