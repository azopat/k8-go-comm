package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewMinioClient - returns new minio client
func NewMinioClient(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) *minio.Client {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return minioClient
}

// UploadFileToMinio - uploads file to minio
func UploadFileToMinio(client *minio.Client, bucketName string, objectName string, reader io.Reader) (minio.UploadInfo, error) {
	uploadInfo, err := client.PutObject(context.Background(), bucketName, objectName, reader, -1, minio.PutObjectOptions{ContentType: "text/plain"})
	if err != nil {
		log.Fatalln(err)
		return uploadInfo, err
	}
	return uploadInfo, nil
}

// GetObjectFromMinio - get file from minio
func GetObjectFromMinio(client *minio.Client, bucketName string, objectName string) (*minio.Object, error) {
	object, err := client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return object, nil
	}
	return object, nil
}

// DeleteObjectInMinio - delete object in minio
func DeleteObjectInMinio(client *minio.Client, bucketName string, objectName string) error {
	err := client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	return err
}

// GetPresignedURLForObject - Get signed url for object
func GetPresignedURLForObject(client *minio.Client, bucketName string, objectName string, expiresIn time.Duration) *url.URL {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	// reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")
	// Generates a presigned url which expires in a day.
	presignedURL, err := client.PresignedGetObject(context.Background(), bucketName, objectName, expiresIn, reqParams)
	if err != nil {
		fmt.Println(err)
		return presignedURL
	}
	return presignedURL
}

func UploadAndReturnURL(client *minio.Client, bucketName string, fileFullPath string, expiresIn time.Duration) *url.URL {

	contentType := getContentType(fileFullPath)
	objectName := filepath.Base(fileFullPath)
	_, err := minioClient.FPutObject(context.Background(), bucketName, objectName, fileFullPath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
		return *url.URL{}, err
	}

	return GetPresignedURLForObject(client, bucketName, objectName, expiresIn)
}

func getContentType(fileFullPath string) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}