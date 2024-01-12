package minioclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client *minio.Client
}

func NewMinioClient() (*MinioClient, error) {

	endpoint := "185.204.2.233:9900"
	accessKey := "Tasks"
	secretKey := "123456789"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &MinioClient{
		Client: minioClient,
	}, nil
}

func GetMinioClient() *minio.Client {

	endpoint := "185.204.2.233:9900"
	accessKey := "Tasks"
	secretKey := "123456789"

	useSsl := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSsl,
	})
	if err != nil {
		panic(err)
	}
	return minioClient
}

// UploadServiceImage загружает изображение в MinIO и возвращает URL изображения.
func (mc *MinioClient) UploadServiceImage(taskID int, imageBytes []byte, contentType string) (string, error) {
	objectName := fmt.Sprintf("tasks/%d/image", taskID)

	// Используйте io.NopCloser вместо ioutil.NopCloser
	reader := io.NopCloser(bytes.NewReader(imageBytes))

	_, err := mc.Client.PutObject(context.TODO(), "images", objectName, reader, int64(len(imageBytes)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	imageURL := fmt.Sprintf("http://185.204.2.233:9900/images/%s", objectName)
	return imageURL, nil
}

// RemoveServiceImage удаляет изображение услуги из MinIO.
func (mc *MinioClient) RemoveServiceImage(taskID int) error {
	objectName := fmt.Sprintf("tasks/%d/image", taskID)
	log.Println(objectName)
	err := mc.Client.RemoveObject(context.TODO(), "images", objectName, minio.RemoveObjectOptions{})
	log.Println(err)
	if err != nil {
		fmt.Println("Failed to remove object from MinIO:", err)
		// Обработка ошибки удаления изображения из MinIO
		return err
	}
	fmt.Println("Image was removed from MinIO successfully:", objectName)
	return nil
}

func ReadObject(bucketName string, objectName string) (contentBytes []byte, contentType string, err error) {
	minioClient := GetMinioClient()
	object, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	contentBytes, err = io.ReadAll(object)
	if err != nil {
		return
	}
	stat, err := object.Stat()
	if err != nil {
		return
	}
	contentType = stat.ContentType
	return
}

// func UploadObject(bucketName string, objectName string, reader io.Reader, size int64, contentType string) error {
// 	minioClient, eror := NewMinioClient()
// 	if eror != nil {
// 		return eror
// 	}
// 	_, err := minioClient.PutObject(context.Background(), bucketName, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func DeleteObject(bucketName string, objectName string) error {
// 	minioClient := GetMinioClient()
// 	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
