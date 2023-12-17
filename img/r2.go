package img

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mime/multipart"
	"net/url"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	cfg "github.com/twibber/api/config"
)

var client *s3.Client

func init() {
	client = R2Client(cfg.Config.R2AccountID, cfg.Config.R2AccessKeyID, cfg.Config.R2AccessKeySecret)
}

func R2Client(accountId, accessKeyId, accessKeySecret string) *s3.Client {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),

		// to prevent "resolve auth scheme: resolve endpoint: endpoint rule error, Invalid region: region was not a valid DNS name."
		config.WithRegion("auto"),
	)
	if err != nil {
		log.WithError(err).Fatal("unable to load AWS config")
	}

	return s3.NewFromConfig(awsCfg)
}

// UploadFile uploads a file to the R2 bucket and returns the file URL
func UploadFile(file *multipart.FileHeader, directory string, filename string) (string, error) {
	// Check if the file is not nil
	if file == nil {
		return "", fmt.Errorf("no file to upload")
	}

	// Check if the file is empty
	if file.Size == 0 {
		return "", fmt.Errorf("file is empty")
	}

	// Ensure the client is initialized
	if client == nil {
		return "", fmt.Errorf("s3 client is not initialized")
	}

	// Validate configuration
	if cfg.Config.R2BucketName == "" {
		return "", fmt.Errorf("bucket name is not configured")
	}

	// Open the file for reading
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	var filePath string
	if directory == "" {
		filePath = filename
	} else {
		filePath = fmt.Sprintf("%s/%s", directory, filename)
	}

	fileExt := filepath.Ext(file.Filename)
	filePath = filePath + fileExt

	// Upload the file
	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &cfg.Config.R2BucketName,
		Key:    &filePath,
		Body:   src,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"bucket": cfg.Config.R2BucketName,
			"key":    filePath,
			"error":  err.Error(),
		}).Error("Failed to upload file to R2 bucket")
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Construct and validate the file URL
	fileURL := fmt.Sprintf("https://cdn.twibber.xyz/%s", filePath)
	if _, err := url.ParseRequestURI(fileURL); err != nil {
		return "", fmt.Errorf("invalid file URL: %s", err)
	}

	return fileURL, nil
}
