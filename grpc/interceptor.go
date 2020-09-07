package grpc

import (
	"context"
	"log"

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
		log.Printf("error: %+v", err)     // スタックトレースを出力
		err = convertErrorWithStatus(err) // ステータス付きのエラーに変換。後述
	}

	// レスポンスを返す
	return m, err
}

func convertErrorWithStatus(err error) error {
	var errorMessage string
	var fieldName string
	var typ string

	// validation エラーの場合
	if _, ok := errors.Cause(err).(validator.ValidationErrors); ok {
		log.Println("ValidationErrors")
		for _, err := range err.(validator.ValidationErrors) {

			fieldName = err.Field()

			switch fieldName {
			case "Content":
				typ = err.Tag()
				switch typ {
				case "max":
					errorMessage = messageContentMax
				case "min":
					errorMessage = messageContentMin
				}
			}

		}
	} else {
		errorMessage = err.Error()
	}

	st := status.New(codes.InvalidArgument, "some error occurred")

	v := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       fieldName,
				Description: errorMessage,
			},
		},
	}
	dt, _ := st.WithDetails(v)

	return dt.Err()
}
