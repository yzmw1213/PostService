package grpc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func transmitStatusInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// メソッドより前に呼ばれる処理

	// メソッドの処理
	m, err := handler(ctx, req)

	// メソッドの処理後に呼ばれる処理
	if err != nil {
		// ステータス付きのエラーに変換。
		err = convertErrorWithStatus(err)
	}

	// レスポンスを返す
	return m, err
}

func convertErrorWithStatus(err error) error {
	var errorStatus string
	var fieldName string
	var typ string

	// validation エラーの場合
	if _, ok := errors.Cause(err).(validator.ValidationErrors); ok {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName = err.Field()
			switch fieldName {
			// 投稿Contentのバリデーションエラー
			case "Content":
				typ = err.Tag()
				switch typ {
				case "max":
					errorStatus = StatusPostContentStringCount
					break
				case "min":
					errorStatus = StatusPostContentStringCount
					break
				}
			// 投稿Titleのバリデーションエラー
			case "Title":
				typ = err.Tag()
				switch typ {
				case "max":
					errorStatus = StatusPostTitleStringCount
					break
				case "min":
					errorStatus = StatusPostTitleStringCount
					break
				}
			// タグ名のバリデーションエラー
			case "TagName":
				typ = err.Tag()
				switch typ {
				case "min":
					errorStatus = StatusTagNameStringCount
					break
				case "max":
					errorStatus = StatusTagNameStringCount
					break
				}
			// コメントContentのバリデーションエラー
			case "CommentContent":
				typ = err.Tag()
				switch typ {
				case "min":
					errorStatus = StatusCommentContentStringCount
					break
				case "max":
					errorStatus = StatusCommentContentStringCount
					break
				}
			}
		}
	} else {
		errorStatus = err.Error()
	}

	st := status.New(codes.InvalidArgument, errorStatus)

	v := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       fieldName,
				Description: errorStatus,
			},
		},
	}
	dt, _ := st.WithDetails(v)

	return dt.Err()
}
