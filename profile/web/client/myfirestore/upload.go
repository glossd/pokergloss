package myfirestore

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/profile/conf"
	"io"
	"mime/multipart"
)

func UploadFile(ctx context.Context, file *multipart.FileHeader, folder, objectName string) (string, error) {
	f, err := file.Open()
	bkt, err := conf.StorageClient.DefaultBucket()
	if err != nil {
		return "", err
	}
	obj := bkt.Object(folder + "/" + objectName)
	wc := obj.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	return buildPublicURI(folder, objectName), nil
}

// https://stackoverflow.com/a/21226442/10160865
func buildPublicURI(folder, objectName string) string {
	return "https://firebasestorage.googleapis.com/v0/b/pocrium.appspot.com/o/" + folder + "%2F" + objectName + "alt=media"
}
