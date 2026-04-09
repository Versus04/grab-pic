package detection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/jackc/pgx/v5"
)

func UserPhoto(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file error", 400)
		return
	}
	defer file.Close()

	Db, err := pgx.Connect(context.TODO(), "postgres://postgres:postgres@localhost:5432/grabpic")
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	// Convert album_id
	albumID, _ := strconv.Atoi(r.FormValue("album_id"))

	// Create temp file
	ext := filepath.Ext(handler.Filename)
	tempFile, _ := os.CreateTemp("temp-images", "upload-*"+ext)
	defer tempFile.Close()

	// 🔥 IMPORTANT: copy file
	io.Copy(tempFile, file)
	tempFile.Seek(0, 0)

	// Decode image
	img, format, err := image.Decode(tempFile)
	if err != nil {
		http.Error(w, "decode error", 500)
		return
	}

	buf := new(bytes.Buffer)
	switch format {
	case "jpeg":
		jpeg.Encode(buf, img, nil)
	case "png":
		png.Encode(buf, img)
	default:
		http.Error(w, "unsupported format", 400)
		return
	}

	// Search in Rekognition
	matchList, err := rekClient.SearchFacesByImage(context.TODO(),
		&rekognition.SearchFacesByImageInput{
			CollectionId:       aws.String("trial-collection"),
			Image:              &types.Image{Bytes: buf.Bytes()},
			FaceMatchThreshold: aws.Float32(90),
			MaxFaces:           aws.Int32(1),
			QualityFilter:      types.QualityFilterAuto,
		})
	if err != nil {
		fmt.Fprint(w, "rekognition error ", err)
		return
	}

	if len(matchList.FaceMatches) == 0 {
		w.Write([]byte("No match found"))
		return
	}

	faceID := matchList.FaceMatches[0].Face.FaceId

	var personID int
	err = Db.QueryRow(context.TODO(),
		`SELECT person_id FROM faces WHERE face_id=$1 LIMIT 1`,
		faceID,
	).Scan(&personID)

	if err != nil {
		http.Error(w, "person not found", 500)
		return
	}

	// Fetch all images
	rows, err := Db.Query(context.TODO(),
		`SELECT i.link
         FROM faces f
         JOIN images i ON f.image_id = i.image_id
         WHERE f.person_id = $1 AND f.album_id = $2`,
		personID, albumID,
	)
	if err != nil {
		http.Error(w, "db query error", 500)
		return
	}
	defer rows.Close()

	var results []string

	for rows.Next() {
		var link string
		rows.Scan(&link)
		results = append(results, link)
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.Encode(results)
}
