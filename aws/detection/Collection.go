package detection

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/joho/godotenv"
)

var rekClient *rekognition.Client

func CreateCollection() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error fetching env file ", err)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		fmt.Println("Error loading configuration ", err)
	}
	collection_id := "trial-collection"
	rekClient = rekognition.NewFromConfig(cfg)
	_, err = rekClient.CreateCollection(context.TODO(), &rekognition.CreateCollectionInput{
		CollectionId: aws.String(collection_id),
	})
	if err != nil {
		fmt.Println("Error creating a new collection ", err)
		return
	}
}
func DeleteCollections(rekClient *rekognition.Client, collection_id string) {
	rekClient.DeleteCollection(context.TODO(), &rekognition.DeleteCollectionInput{CollectionId: aws.String(collection_id)})
}
func ListCollections() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error fetching env file ", err)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading configuration ", err)
	}
	rekClient = rekognition.NewFromConfig(cfg)
	collectionList, err := rekClient.ListCollections(context.TODO(), &rekognition.ListCollectionsInput{})
	if err != nil {
		fmt.Println("Error listing collections ", err)
	}
	list := collectionList.CollectionIds
	for _, b := range list {
		DeleteCollections(rekClient, b)
		fmt.Println("Deleted Collection ", b)
	}
}
