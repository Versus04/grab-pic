package main

import (
	"net/http"

	"github.com/Versus04/grab-pic/FileUploader"
	"github.com/Versus04/grab-pic/aws/detection"
	_ "github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	detection.CreateCollection()
	http.HandleFunc("/upload", FileUploader.FileUploader)

	http.HandleFunc("/getphoto", detection.UserPhoto)
	http.ListenAndServe("localhost:8080", nil)

}
