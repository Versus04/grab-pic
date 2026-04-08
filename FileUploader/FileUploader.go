package FileUploader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Versus04/grab-pic/aws"
	"github.com/Versus04/grab-pic/aws/detection"
	"github.com/jackc/pgx/v5"
)

var Db *pgx.Conn

func FileUploader(w http.ResponseWriter, r *http.Request) {
	detection.CreateCollection()
	r.ParseMultipartForm(20 << 20)
	Db, _ = pgx.Connect(context.TODO(), "postgres://postgres:postgres@localhost:5432/grabpic")

	fileList := r.MultipartForm.File["files"]
	name := "goa trip"
	var album_id int
	err := Db.QueryRow(context.TODO(), `INSERT INTO albums(name) VALUES ($1) RETURNING album_id`, name).Scan(&album_id)
	if err != nil {
		fmt.Println("Error executing database")
		return
	}
	for _, fh := range fileList {
		ext := filepath.Ext(fh.Filename)
		file, err := fh.Open()
		temp, err := os.CreateTemp("temp-images", "upload-*"+ext)
		if err != nil {
			return
		}
		_, err = io.Copy(temp, file)
		url, err := aws.UploadToS3(temp)
		if err != nil {
			fmt.Println("Error putting or getting object ", err)
			return
		}

		var image_id int
		err = Db.QueryRow(context.TODO(), `INSERT INTO images(album_id,link) VALUES ($1,$2) RETURNING image_id`, album_id, url).Scan(&image_id)

		detection.Detect("uploads/"+filepath.Base(temp.Name()), url, temp, image_id, album_id, Db)
		temp.Close()
		file.Close()

	}
}
