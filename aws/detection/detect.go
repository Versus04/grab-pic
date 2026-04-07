package detection

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image"
	"image/jpeg"
	_ "image/jpeg"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/jackc/pgx/v5"
)

func CreateImageBuffer(file *os.File, imageHeight int, imageWidth int, bbox *types.BoundingBox) (bytes.Buffer, error) {
	return bytes.Buffer{}, nil
}
func Detect(filename string, link string, file *os.File, image_id int, album_id int, Db *pgx.Conn) {
	collection_id := "trial-collection"
	images := &types.Image{

		S3Object: &types.S3Object{
			Bucket: aws.String("photo-bucket-67"),
			Name:   aws.String(filename),
		},
	}
	file.Seek(0, 0)
	detectedFaces, err := rekClient.DetectFaces(context.TODO(), &rekognition.DetectFacesInput{
		Image: images,
	})

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("Error decoding config:", err)
		return
	}
	file.Seek(0, 0)
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image ", err)
		return
	}
	file.Seek(0, 0)
	imageHeight := config.Height
	imageWidth := config.Width
	for _, b := range detectedFaces.FaceDetails {
		bbox := b.BoundingBox
		x := int((*bbox.Left) * float32(imageWidth))
		y := int((*bbox.Top) * float32(imageHeight))
		w := int((*bbox.Width) * float32(imageWidth))
		h := int((*bbox.Height) * float32(imageHeight))
		rect := image.Rect(x, y, x+w, y+h)
		faceImg := img.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(rect)
		buf := new(bytes.Buffer)
		err := jpeg.Encode(buf, faceImg, nil)
		if err != nil {
			fmt.Println("Error in the buffer ", err)
			return
		}
		input := &rekognition.SearchFacesByImageInput{
			CollectionId:       aws.String(collection_id),
			Image:              &types.Image{Bytes: buf.Bytes()},
			FaceMatchThreshold: aws.Float32(90),
			MaxFaces:           aws.Int32(1),
			QualityFilter:      types.QualityFilterAuto,
		}
		searchFaceList, err := rekClient.SearchFacesByImage(context.TODO(), input)
		if err != nil {
			fmt.Println("Error connecting to aws Rekognition ", err)
			return
		}
		if len(searchFaceList.FaceMatches) > 0 {

			var person_id int
			faceID := searchFaceList.FaceMatches[0].Face.FaceId
			_ = Db.QueryRow(context.TODO(), `SELECT person_id FROM faces WHERE face_id=$1`, faceID).Scan(&person_id)
			_, err = Db.Exec(context.TODO(), `INSERT INTO faces(face_id,person_id,image_id,album_id) VALUES($1,$2,$3,$4)`, faceID, person_id, image_id, album_id)
		} else {
			input := &rekognition.IndexFacesInput{
				CollectionId: aws.String(collection_id),
				Image:        &types.Image{Bytes: buf.Bytes()},
			}
			indexFaceList, err := rekClient.IndexFaces(context.TODO(), input)
			if err != nil {
				fmt.Println("Error indexing Face ", err)
			}
			var personId int
			err = Db.QueryRow(context.TODO(), `INSERT INTO persons DEFAULT VALUES RETURNING person_id`).Scan(&personId)
			if err != nil {
				fmt.Println("error getting person_id ", err)
				return
			}
			if len(indexFaceList.FaceRecords) == 0 {
				fmt.Println("No face indexed")
				continue
			}
			faceId := indexFaceList.FaceRecords[0].Face.FaceId
			_, err = Db.Exec(context.TODO(), `INSERT INTO faces(face_id,person_id,image_id,album_id) VALUES ($1,$2,$3,$4)`, faceId, personId, image_id, album_id)
			if err != nil {
				fmt.Println("Error inserting into face database ", err)
				return
			}
		}
	}
}
