// https://godoc.org/cloud.google.com/go/storage

package persistence

import (
	"context"
	"io"
	"log"

	"cloud.google.com/go/storage"
)

// GoogleCloudStorage represents a struct for a GCS client
type GoogleCloudStorage struct {
	bucketName string
	objectName string
	context    context.Context
	client     *storage.Client
}

// NewGoogleCloudStorage returns an initiated *GoogleCloudStorage
// returns nil on error while Client initialization
func NewGoogleCloudStorage(bucketName, objectName string) *GoogleCloudStorage {
	context := context.Background()
	client, err := storage.NewClient(context)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return nil
	}

	return &GoogleCloudStorage{
		context:    context,
		bucketName: bucketName,
		objectName: objectName,
		client:     client,
	}
}

// NewReader returns a new io.Reader for the GoogleCloudStorage
func (g *GoogleCloudStorage) NewReader() (r io.ReadCloser, err error) {
	bucket := g.client.Bucket(g.bucketName)
	object := bucket.Object(g.objectName)
	return object.NewReader(g.context)
}

// NewWriter returns a new io.Reader for the GoogleCloudStorage
func (g *GoogleCloudStorage) NewWriter() io.WriteCloser {
	bucket := g.client.Bucket(g.bucketName)
	object := bucket.Object(g.objectName)
	return object.NewWriter(g.context)
}
