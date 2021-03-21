package aws

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	sess               *session.Session
	awsAccessKey       = os.Getenv("AWS_ACCESS_KEY")
	s3_bucket          = os.Getenv("AWS_S3_BUCKET_NAME")
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint           = os.Getenv("AWS_S3_ENDPOINT")
	region             = os.Getenv("AWS_S3_REGION")
	downloadDir        = os.Getenv("OBJECT_DOWNLOAD_DIR")
	letters            = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

type downloader struct {
	bucket string
	file   string
	dir    string
	*s3manager.Downloader
}

//GetS3Session gets session for Connecting to S3
func GetS3Session() {
	session, err := NewSession()
	if err != nil {
		fmt.Println(err)
	}
	sess = session
}

func Upload(imageBase64 string) (string, error) {
	GetS3Session()
	// ファイルを開く
	DATE := time.Now().Format("2006-01-02")
	NAME := randSeq(15)
	key := fmt.Sprintf("%s/%s", DATE, NAME)

	uploader := s3manager.NewUploader(sess)

	data, _ := base64.StdEncoding.DecodeString(imageBase64)
	wb := new(bytes.Buffer)
	wb.Write(data)

	uo, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &s3_bucket,
		Key:    &key,
		Body:   wb,
	})

	log.Println("bucket", s3_bucket)
	log.Println("key", key)
	log.Println("location", uo.Location)
	if err != nil {
		log.Println(err)
	}
	S3_END := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", s3_bucket, region)

	return strings.Replace(uo.Location, S3_END, "", 1), err
}

// randSeq 指定した文字数のランダム文字列を返却
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
