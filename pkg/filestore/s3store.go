package filestore

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime/multipart"
	"os"
	"schooli-api/internal/models"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

type S3Storage struct {
	client     *minio.Client
	endpoint   string
	accessKey  string
	secretKey  string
	secure     bool
	signV2     bool
	region     string
	bucket     string
	pathPrefix string
	encrypt    bool
	trace      bool
	skipVerify bool
	timeout    time.Duration
}

func NewS3Storage(settings FileBackendSettings) (*S3Storage, error) {
	timeout := time.Duration(settings.AmazonS3RequestTimeoutMilliseconds) * time.Millisecond
	backend := &S3Storage{
		endpoint:   settings.AmazonS3Endpoint,
		accessKey:  settings.AmazonS3AccessKeyId,
		secretKey:  settings.AmazonS3SecretAccessKey,
		secure:     settings.AmazonS3SSL,
		signV2:     settings.AmazonS3SignV2,
		region:     settings.AmazonS3Region,
		bucket:     settings.AmazonS3Bucket,
		pathPrefix: settings.AmazonS3PathPrefix,
		encrypt:    settings.AmazonS3SSE,
		trace:      settings.AmazonS3Trace,
		skipVerify: settings.SkipVerify,
		timeout:    timeout,
	}
	cli, err := backend.s3New()
	if err != nil {
		return nil, err
	}
	backend.client = cli
	slog.Info("S3 client connected")
	return backend, nil
}

// Similar to s3.New() but allows initialization of signature v2 or signature v4 client.
// If signV2 input is false, function always returns signature v4.
//
// Additionally, this function also takes a user defined region, if set
// disables automatic region lookup.
func (s *S3Storage) s3New() (*minio.Client, error) {
	opts := minio.Options{
		Creds:  credentials.NewStaticV4(s.accessKey, s.secretKey, ""),
		Secure: s.secure,
		Region: s.region,
	}

	tr, err := minio.DefaultTransport(s.secure)
	if err != nil {
		return nil, err
	}
	if s.skipVerify {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	opts.Transport = tr

	// If this is a cloud installation, we override the default transport.
	//TODO override in cloud
	//if isCloud {
	//	scheme := "http"
	//	if s.secure {
	//		scheme = "https"
	//	}
	//	newTransport := http.DefaultTransport.(*http.Transport).Clone()
	//	newTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: s.skipVerify}
	//	opts.Transport = &customTransport{
	//		host:   s.endpoint,
	//		scheme: scheme,
	//		client: http.Client{Transport: newTransport},
	//	}
	//}

	s3client, err := minio.New(s.endpoint, &opts)
	if err != nil {
		return nil, err
	}

	if s.trace {
		s3client.TraceOn(os.Stdout)
	}

	return s3client, nil
}

func (s *S3Storage) UploadFile(ctx context.Context, fum models.FileUploadModel) (minio.UploadInfo, error) {
	folder := fmt.Sprintf("%v_%s", fum.FolderName, fum.ObjName)
	info, err := s.client.PutObject(
		ctx,
		s.bucket,
		folder,
		fum.FileBuf,
		fum.FileSize,
		minio.PutObjectOptions{ContentType: fum.ContentType})
	if err != nil {
		return info, err
	}
	return info, nil
}
func (s *S3Storage) MultipleFileUpload(ctx context.Context, mf *models.MultipleFileUploadModel) ([]string, []error) {
	type item struct {
		key string
		err error
	}

	ch := make(chan item, len(mf.FileNames))

	for _, f := range mf.FileNames {
		go func(ctx context.Context, fh *multipart.FileHeader) {
			ctt := fh.Header["Content-Type"][0]
			v, _ := fh.Open()
			var it item
			filename := fmt.Sprintf("%s/%s", mf.FolderName, fh.Filename)
			info, err := s.client.PutObject(ctx, s.bucket, filename, v, fh.Size, minio.PutObjectOptions{ContentType: ctt})
			if err != nil {
				it.err = err
			} else {
				it.key = info.Key
			}
			ch <- it
		}(ctx, f)
	}

	keys := make([]string, 0, len(mf.FileNames))
	errs := make([]error, 0, len(mf.FileNames))

	for range mf.FileNames {
		it := <-ch
		if it.err != nil {
			errs = append(errs, it.err)
		} else {
			keys = append(keys, it.key)
		}
	}

	return keys, errs
}
func (s *S3Storage) ListFiles(ctx context.Context, folderPrefix string) ([]string, error) {
	newctx, cancel := context.WithCancel(ctx)
	defer cancel()

	objectCh := s.client.ListObjects(newctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    folderPrefix,
		Recursive: true,
	})
	imgs := make([]string, len(objectCh))

	for object := range objectCh {
		if object.Err != nil {
			slog.Error("error listing directory", object.Err)
			return nil, object.Err
		}
		f := fmt.Sprintf("%s/%s/%s", s.client.EndpointURL(), s.bucket, object.Key)
		imgs = append(imgs, f)
	}
	return imgs, nil
}
func (s *S3Storage) DeleteFolder(ctx context.Context, folderPrefix string) error {
	newctx, cancel := context.WithCancel(ctx)

	defer cancel()

	objectsCh := s.client.ListObjects(newctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    folderPrefix,
		Recursive: true,
	})

	for object := range objectsCh {
		if object.Err != nil {
			slog.Error("error listing directory", object.Err)
			return object.Err
		}
		err := s.client.RemoveObject(newctx, s.bucket, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *S3Storage) DeleteFolderFile(ctx context.Context, folderPrefix string) error {
	err := s.client.RemoveObject(ctx, s.bucket, folderPrefix, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
func (s *S3Storage) MakeBucket() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{Region: s.region})
	if err != nil {
		return errors.Wrap(err, "unable to create the s3 bucket")
	}
	return nil
}
func (s *S3Storage) BuildFilePath(path string) string {
	return fmt.Sprintf("%s/%s/%s", s.client.EndpointURL(), s.bucket, path)
}
