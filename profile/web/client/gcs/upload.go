package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	goauth "github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/goconf"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"time"
)

var client *storage.Client

func Init() context.CancelFunc {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	var err error
	client, err = storage.NewClient(ctx, goauth.GoogleClientOptions()...)
	if err != nil {
		log.Errorf("Couldn't create google storage client: %v", err)
	}
	return cancel
}

func UploadFile(ctx context.Context, file *multipart.FileHeader, bucketName, objectName string) (string, error) {
	if client == nil {
		log.Errorf("User tried to upload avatar, but GCS client is not initialized")
		return "", fmt.Errorf("avatar storage is not available, try later")
	}

	f, err := file.Open()
	if err != nil {
		return "", err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(objectName)
	wc := obj.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	return buildPublicURI(objectName), nil
}

// https://stackoverflow.com/a/21226442/10160865
func buildPublicURI(objectName string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", goconf.Props.AvatarBucket, objectName)
}
