package minio

import (
	"context"
	minioClient "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/encrypt"
	"io"
)

type Client struct {
	ctx             context.Context
	endpoint        string
	region          string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	password        string
	client          *minioClient.Client
}

func New(opts ...Option) (*Client, error) {
	c := &Client{
		ctx:      context.Background(),
		endpoint: "127.0.0.1:9000",
		region:   "cn-north-1",
	}
	for _, opt := range opts {
		opt(c)
	}

	// Initialize minio client object.
	minioClient, err := minioClient.New(c.endpoint, &minioClient.Options{
		Creds:  credentials.NewStaticV4(c.accessKeyID, c.secretAccessKey, ""),
		Secure: c.useSSL,
	})
	if err != nil {
		return nil, err
	}
	c.client = minioClient
	return c, nil
}

// MakeBucket Make Bucket
func (c *Client) MakeBucket(bucketName string, objectLocking ...bool) error {
	locking := false
	if len(objectLocking) > 0 {
		locking = objectLocking[0]
	}
	err := c.client.MakeBucket(c.ctx, bucketName, minioClient.MakeBucketOptions{Region: c.region, ObjectLocking: locking})
	if err != nil {
		//Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := c.client.BucketExists(c.ctx, bucketName)
		if errBucketExists == nil && exists {
		} else {
			return err
		}
	}
	return nil
}

// PutObject Put Object
func (c *Client) PutObject(bucketName string, objectName string, reader io.Reader, objectSize int64, userTags map[string]string) (*UploadInfo, error) {
	opts := minioClient.PutObjectOptions{}
	if userTags != nil {
		opts.UserTags = userTags
		opts.UserMetadata = userTags
	}

	if c.password != "" {
		opts.ServerSideEncryption = encrypt.DefaultPBKDF([]byte(c.password), []byte(bucketName+objectName))
	}

	info, err := c.client.PutObject(c.ctx, bucketName, objectName, reader, objectSize, opts)
	if err != nil {
		return nil, err
	}
	return &UploadInfo{
		Bucket:       info.Bucket,
		Key:          info.Key,
		ETag:         info.ETag,
		Size:         info.Size,
		LastModified: info.LastModified,
	}, nil
}

// FPutObject Put Object
func (c *Client) FPutObject(bucketName string, objectName string, filePath string, userTags map[string]string) (*UploadInfo, error) {
	opts := minioClient.PutObjectOptions{}
	if userTags != nil {
		opts.UserTags = userTags
		opts.UserMetadata = userTags
	}

	if c.password != "" {
		opts.ServerSideEncryption = encrypt.DefaultPBKDF([]byte(c.password), []byte(bucketName+objectName))
	}

	info, err := c.client.FPutObject(c.ctx, bucketName, objectName, filePath, opts)
	if err != nil {
		return nil, err
	}

	return &UploadInfo{
		Bucket:       info.Bucket,
		Key:          info.Key,
		ETag:         info.ETag,
		Size:         info.Size,
		LastModified: info.LastModified,
	}, nil
}

// GetObjectByte Get Object Byte
func (c *Client) GetObjectByte(bucketName, objectName string) ([]byte, error) {
	opts := minioClient.GetObjectOptions{}

	if c.password != "" {
		opts.ServerSideEncryption = encrypt.DefaultPBKDF([]byte(c.password), []byte(bucketName+objectName))
	}

	object, err := c.client.GetObject(c.ctx, bucketName, objectName, opts)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// GetObject Get Object
func (c *Client) GetObject(bucketName, objectName string) (*minioClient.Object, error) {
	opts := minioClient.GetObjectOptions{}

	if c.password != "" {
		opts.ServerSideEncryption = encrypt.DefaultPBKDF([]byte(c.password), []byte(bucketName+objectName))
	}

	object, err := c.client.GetObject(c.ctx, bucketName, objectName, opts)
	if err != nil {
		return nil, err
	}
	return object, nil
}

//
//// GetObjectRead Get Object read
//func (c *Client) GetObjectRead(bucketName, objectName string) (io.ReadCloser, *minioClient.ObjectInfo, *http.Header, error) {
//	opts := minioClient.GetObjectOptions{}
//
//	if c.password != "" {
//		opts.ServerSideEncryption = encrypt.DefaultPBKDF([]byte(c.password), []byte(bucketName+objectName))
//	}
//
//	readCloser, objectInfo, header, err := c.client.GetObjectRead(c.ctx, bucketName, objectName, opts)
//
//	return readCloser, &objectInfo, &header, err
//}

// StatObject Stat Object
func (c *Client) StatObject(bucketName, objectName string) (*ObjectInfo, error) {
	info, err := c.client.StatObject(c.ctx, bucketName, objectName, minioClient.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &ObjectInfo{
		ETag:           info.ETag,
		Key:            info.Key,
		LastModified:   info.LastModified,
		Size:           info.Size,
		ContentType:    info.ContentType,
		UserMetadata:   info.UserMetadata,
		UserTags:       info.UserTags,
		UserTagCount:   info.UserTagCount,
		StorageClass:   info.StorageClass,
		IsLatest:       info.IsLatest,
		IsDeleteMarker: info.IsDeleteMarker,
		VersionID:      info.VersionID,
	}, nil
}

// RemoveObject Remove Object
func (c *Client) RemoveObject(bucketName, objectName string) error {
	return c.client.RemoveObject(c.ctx, bucketName, objectName, minioClient.RemoveObjectOptions{})
}
