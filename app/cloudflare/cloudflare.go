package cloudflare

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// cloudflareに画像をアップロードします
//
// 公開URLを返します。
func Upload(c *gin.Context, image *multipart.FileHeader) (string, error) {
	var bucketName = "gacha-bot"
	var accessKeyId = os.Getenv("CLOUDFLARE_S3_ACCESS_KEY")
	var accessKeySecret = os.Getenv("CLOUDFLARE_S3_SECRET_ACCESS_KEY")

	// 1. *multipart.FileHeaderからmultipart.Fileをオープンする
	imageFile, err := image.Open()
	if err != nil {
		return "", errors.NewError("画像のオープンに失敗しました", err)
	}
	defer imageFile.Close()

	// uuidを生成します
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", errors.NewError("UUIDの生成に失敗しました", err)
	}

	publicURL := fmt.Sprintf(
		"%s/%s",
		os.Getenv("CLOUDFLARE_URL"),
		newUUID.String(),
	)

	var f = func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{URL: os.Getenv("CLOUDFLARE_S3_API_URL")}, nil
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(f)

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKeyId, accessKeySecret, "",
			),
		),
	)
	if err != nil {
		return "", errors.NewError("設定の読み込みに失敗しました", err)
	}

	// MIMEタイプの検出
	buffer := make([]byte, 512) // ヘッダの最大長
	_, err = imageFile.Read(buffer)
	if err != nil {
		return "", errors.NewError("画像の読み込みに失敗しました", err)
	}
	contentType := http.DetectContentType(buffer)
	_, err = imageFile.Seek(0, 0) // ファイルのポインタを先頭に戻す
	if err != nil {
		return "", errors.NewError("MIMEタイプの検出に失敗しました", err)
	}

	client := s3.NewFromConfig(cfg)
	param := s3.PutObjectInput{
		Bucket:             aws.String(bucketName),
		Key:                aws.String(newUUID.String()),
		Body:               imageFile,
		ContentType:        aws.String(contentType),
		ContentDisposition: aws.String("inline"),
	}
	_, err = client.PutObject(context.Background(), &param)
	if err != nil {
		return "", errors.NewError("画像のアップロードに失敗しました", err)
	}

	return publicURL, nil
}
