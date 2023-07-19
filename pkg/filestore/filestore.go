package filestore

import (
	"context"
	"schooli-api/internal/models"

	"github.com/minio/minio-go/v7"
)

type FileStorage interface {
	UploadFile(ctx context.Context, fum models.FileUploadModel) (minio.UploadInfo, error)
	MultipleFileUpload(ctx context.Context, mf *models.MultipleFileUploadModel) ([]string, []error)
	ListFiles(ctx context.Context, folderPrefix string) ([]string, error)
	DeleteFolder(ctx context.Context, folderPrefix string) error
	DeleteFolderFile(ctx context.Context, folderPrefix string) error
	BuildFilePath(path string) string
}

type FileBackendSettings struct {
	AmazonS3AccessKeyId                string
	AmazonS3SecretAccessKey            string
	AmazonS3Bucket                     string
	AmazonS3PathPrefix                 string
	AmazonS3Region                     string
	AmazonS3Endpoint                   string
	AmazonS3SSL                        bool
	AmazonS3SignV2                     bool
	AmazonS3SSE                        bool
	AmazonS3Trace                      bool
	SkipVerify                         bool
	AmazonS3RequestTimeoutMilliseconds int64
}
