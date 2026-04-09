package main

import (
	"github.com/Versus04/grab-pic/aws/detection"
	_ "github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	detection.ListCollections()
	//mux := http.NewServeMux()
	//detection.CreateCollection()
	//mux.HandleFunc("/upload", FileUploader.FileUploader)

	//mux.HandleFunc("/getphoto", detection.UserPhoto)
	//http.ListenAndServe(":8080", corsMiddleware(mux))

}
