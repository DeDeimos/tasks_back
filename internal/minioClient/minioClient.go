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

	endpoint := "localhost:9000"
	accessKey := "UT3TiBIF0qu15Nel3XcJ"
	secretKey := "rvGlJ2OfR7HQYMIKQsmC7tNoyetGRYUl3TrQ0Mnt"

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

	imageURL := fmt.Sprintf("http://localhost:9000/images/%s", objectName)
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
