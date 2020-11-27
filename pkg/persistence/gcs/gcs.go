package gcs

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

// https://godoc.org/cloud.google.com/go/storage

// GoogleCloudStorage represents a struct for a GCS client
type GoogleCloudStorage struct {
	context    context.Context
	bucketName string
	objectName string
	client     *storage.Client
}

// NewGoogleCloudStorage returns an initiated *GoogleCloudStorage
// returns nil on error while Client initialization
func NewGoogleCloudStorage(ctx context.Context, bucketName, objectName string) (*GoogleCloudStorage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GoogleCloudStorage{
		context:    ctx,
		bucketName: bucketName,
		objectName: objectName,
		client:     client,
	}, nil
}

// NewReader returns a new io.Reader for the GoogleCloudStorage
func (g *GoogleCloudStorage) NewReader() (io.ReadCloser, error) {
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
