package main

import (
	"net/http"

	"github.com/Versus04/grab-pic/aws/detection"
	_ "github.com/aws/aws-sdk-go-v2/config"
)

func main() {

	//http.HandleFunc("/upload", FileUploader.FileUploader)
	detection.CreateCollection()
	http.ListenAndServe("localhost:8080", nil)

}
