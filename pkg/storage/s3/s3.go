package s3

import (
	"context"
	"io"
)

// S3 stands for AWS storage
type S3 struct {
	bucketName string
}

// New returns new S3 object
func New(bucketName string) (*S3, error) {
	return &S3{
		bucketName: bucketName,
	}, nil
}

// Write object to S3
func (s *S3) Write(ctx context.Context, r io.Reader, objName string, public bool) error {
	return nil
}

// Delete object from S3
func (s *S3) Delete(ctx context.Context, objName string) error {
	return nil
}

// Close storage
func (s *S3) Close() error {
	return nil
}
