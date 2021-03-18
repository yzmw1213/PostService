package aws

import (
	"encoding/base64"
	"os"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestS3Session(t *testing.T) {
	GetS3Session()
	if sess == nil {
		t.Fatal("AWS SDKからセッションを取得できませんでした")
	}
}

func TestUpload(t *testing.T) {
	file, _ := os.Open("example.jpeg")
	defer file.Close()

	fi, _ := file.Stat()
	size := fi.Size()

	data := make([]byte, size)
	file.Read(data)

	imgBase64 := base64.StdEncoding.EncodeToString(data)
	location, err := Upload(imgBase64)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", location)
}
