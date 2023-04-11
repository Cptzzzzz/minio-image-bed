package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var UploadUrl = "http://127.0.0.1:8080/"

func main() {
	content, _ := os.ReadFile(os.Args[1])
	buf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(buf)
	w, _ := bodyWriter.CreateFormFile("file", os.Args[1])
	w.Write(content)
	bodyWriter.Close()
	resp, err := http.Post(UploadUrl, bodyWriter.FormDataContentType(), buf)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data))
}
