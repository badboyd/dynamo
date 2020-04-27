package gcs

import (
	"context"
	"hash/crc32"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
)

// Transport stands for google cloud storage transport
type Transport struct {
	client *storage.Client
}

// NewTransport returns new GCSTransport
func NewTransport() http.RoundTripper {
	client, err := storage.NewClient(context.Background())

	if err != nil {
		log.Fatalf("Can't create GCS client: %s", err.Error())
	}

	return Transport{client}
}

// RoundTrip for GCS transport
func (t Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	switch req.Method {
	case http.MethodPut:
		return t.Write(req)
	case http.MethodDelete:
		return t.Delete(req)
	default:
		// should never happens
		return nil, err
	}
}

// Write object to GCS
func (t Transport) Write(req *http.Request) (resp *http.Response, err error) {
	bkt := t.client.Bucket(req.URL.Host)
	obj := bkt.Object(strings.TrimPrefix(req.URL.Path, "/"))

	data, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	ow := obj.NewWriter(context.Background())

	ow.CRC32C = crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))
	ow.SendCRC32C = true

	if _, err = ow.Write(data); err != nil {
		return nil, err
	}

	if err := ow.Close(); err != nil {
		return nil, err
	}

	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     make(http.Header),
		Close:      true,
		Request:    req,
	}, nil
}

// Delete object from GCS
func (t Transport) Delete(req *http.Request) (resp *http.Response, err error) {
	bkt := t.client.Bucket(req.URL.Host)
	obj := bkt.Object(strings.TrimPrefix(req.URL.Path, "/"))
	if err := obj.Delete(context.Background()); err != nil {
		return nil, err
	}

	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     make(http.Header),
		Close:      true,
		Request:    req,
	}, nil
}
