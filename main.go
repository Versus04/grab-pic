package main

import (
	"net/http"

	"github.com/Versus04/grab-pic/FileUploader"
	"github.com/Versus04/grab-pic/aws/detection"
	_ "github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	//detection.ListCollections()
	mux := http.NewServeMux()
	//detection.CreateCollection()
	mux.HandleFunc("/upload", FileUploader.FileUploader)

	mux.HandleFunc("/getphoto", detection.UserPhoto)
	loggedMux := detection.UsageLogger(FileUploader.Db, mux)
	finalHandler := corsMiddleware(loggedMux)

	http.ListenAndServe(":8080", finalHandler)

}
