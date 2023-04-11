package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

var minioClient *minio.Client

var bucketName string
var endpoint string

func prepare() {
	minioClient = nil
	endpoint = os.Args[1]
	accessKeyID := os.Args[2]
	secretAccessKey := os.Args[3]
	bucketName = os.Args[4]
	useSSL := false
	minio, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println(err)
	}
	minioClient = minio
}

func main() {
	prepare()
	http.HandleFunc("/", upload)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func generateFileName(name string) string {
	return fmt.Sprintf("%v%s", time.Now().Unix(), name)
}

func upload(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, _ := r.FormFile("file")
	content, _ := io.ReadAll(file)
	filename := generateFileName(path.Ext(fileHeader.Filename))
	info, err := minioClient.PutObject(context.Background(), bucketName, filename,
		bytes.NewReader(content), fileHeader.Size, minio.PutObjectOptions{
			ContentType: http.DetectContentType(content),
		})
	fmt.Println(info)
	if err != nil {
		fmt.Println(err)
	} else {
		w.Write([]byte(fmt.Sprintf("%s/%s/%s", endpoint, bucketName, filename)))
	}
}
