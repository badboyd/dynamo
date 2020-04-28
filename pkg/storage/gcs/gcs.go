package gcs

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

// GCS stands for google cloud storage
type GCS struct {
	client     *storage.Client
	bucketName string
}

// New returns new GCS
func New(bucketName string) (*GCS, error) {
	client, err := storage.NewClient(context.Background())

	if err != nil {
		return nil, err
	}

	return &GCS{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// Write object to GCS
func (g *GCS) Write(ctx context.Context, r io.Reader, objName string, public bool) error {
	bkt := g.client.Bucket(g.bucketName)
	obj := bkt.Object(objName)

	ow := obj.NewWriter(ctx)

	if _, err := io.Copy(ow, r); err != nil {
		return err
	}

	if err := ow.Close(); err != nil {
		return err
	}

	if public {
		if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return err
		}
	}

	return nil
}

// Delete object from GCS
func (g *GCS) Delete(ctx context.Context, objName string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	bkt := g.client.Bucket(g.bucketName)
	obj := bkt.Object(objName)

	return obj.Delete(ctx)
}

// Close storage
func (g *GCS) Close() error {
	return g.client.Close()
}
