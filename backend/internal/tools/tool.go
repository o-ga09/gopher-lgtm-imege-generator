package tools

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

type GenerateImageInput struct {
	Prompt   string `json:"prompt"`
	Filename string `json:"filename"`
}

type GenerateImageResult struct {
	Filename string `json:"filename"`
	Status   string `json:"Status"`
}

type SaveImageInput struct {
	Filename string `json:"filename"`
}

type SaveImageResult struct {
	Status string `json:"Status"`
}

func GenerateImage(ctx tool.Context, input GenerateImageInput) (GenerateImageResult, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return GenerateImageResult{
			Filename: "",
			Status:   "failed to create genai client",
		}, err
	}
	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-image",
		genai.Text(input.Prompt),
		nil,
	)

	// Create an "tmp" directory in the current working directory if it doesn't exist.
	outputDir := "tmp"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Printf("Failed to create tmp directory '%s': %v", outputDir, err)
		return GenerateImageResult{
			Filename: "",
			Status:   "failed to create tmp directory",
		}, err
	}

	var imageBytes []byte
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Println(part.Text)
		} else if part.InlineData != nil {
			imageBytes = part.InlineData.Data
			_ = os.WriteFile(filepath.Join(outputDir, input.Filename), imageBytes, 0644)
		}
	}

	_, err = ctx.Artifacts().Save(ctx, input.Filename, genai.NewPartFromBytes(imageBytes, "image/png"))
	if err != nil {
		return GenerateImageResult{
			Filename: "",
			Status:   "failed to save artifact",
		}, err
	}
	return GenerateImageResult{Filename: input.Filename, Status: "success"}, nil
}

// createS3Client creates an S3-compatible client for Cloudflare R2 or MinIO
func createS3Client(ctx context.Context) (*s3.Client, error) {
	// Get configuration from environment variables
	accessKey := os.Getenv("CLOUDFLARE_R2_ACCESSKEY")
	secretKey := os.Getenv("CLOUDFLARE_R2_SECRETKEY")
	endpoint := os.Getenv("CLOUDFLARE_R2_ENDPOINT")
	region := os.Getenv("CLOUDFLARE_R2_REGION")
	env := os.Getenv("ENV")

	// Set defaults
	if region == "" {
		region = "auto"
	}
	if env == "local" && endpoint == "" {
		endpoint = "http://minio:9000"
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           endpoint,
			SigningRegion: "auto",
			Source:        aws.EndpointSourceCustom,
		}, nil
	})

	// Create AWS config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithDefaultRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}
	// Create S3 client with custom endpoint
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Required for MinIO and some S3-compatible services
	})

	return client, nil
}

func SaveImage(ctx tool.Context, input SaveImageInput) (SaveImageResult, error) {
	filename := input.Filename
	resp, err := ctx.Artifacts().Load(ctx, filename)
	if err != nil {
		log.Printf("Failed to load artifact '%s': %v", filename, err)
		return SaveImageResult{
			Status: "failed to load artifact",
		}, err
	}

	if resp.Part.InlineData == nil || len(resp.Part.InlineData.Data) == 0 {
		log.Printf("Artifact '%s' has no inline data", filename)
		return SaveImageResult{
			Status: "artifact has no inline data",
		}, err
	}

	// Ensure the filename has a .png extension
	fileName := filename
	if filepath.Ext(fileName) != ".png" {
		fileName += ".png"
	}

	// Generate unique filename with UUID
	uniqueFileName := fmt.Sprintf("%s-%s", uuid.New().String(), fileName)

	// Create S3 client
	s3Client, err := createS3Client(ctx)
	if err != nil {
		log.Printf("Failed to create S3 client: %v", err)
		return SaveImageResult{
			Status: "failed to create S3 client",
		}, err
	}

	// Upload to S3-compatible storage (R2/MinIO)
	bucketName := os.Getenv("CLOUDFLARE_R2_BUCKET_NAME")
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(uniqueFileName),
		Body:        bytes.NewReader(resp.Part.InlineData.Data),
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		log.Printf("Failed to upload image to S3: %v", err)
		return SaveImageResult{
			Status: "failed to upload image to S3",
		}, err
	}

	log.Printf("Successfully uploaded image to S3: %s", uniqueFileName)
	return SaveImageResult{Status: fmt.Sprintf("success: uploaded as %s", uniqueFileName)}, nil
}
