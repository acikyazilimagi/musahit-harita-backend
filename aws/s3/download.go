package s3

import (
	"fmt"
	log2 "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type ObjectData []byte

func (o ObjectData) Bytes() []byte {
	return o
}

func Download(bucket string, key string) ObjectData {
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)

	var data []byte
	buff := &aws.WriteAtBuffer{}
	_, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		fmt.Println(err)
		return data
	}
	return buff.Bytes()
}

func DownloadMostRecentObject(bucket string) ObjectData {
	contents, err := ListObjects(bucket)

	if err != nil {
		log2.Logger().Error("Couldn't list objects in bucket %v. Here's why: %v\n", zap.String("bucket", bucket), zap.Error(err))
		return nil
	}
	if len(contents) == 0 {
		log2.Logger().Warn("No objects in bucket %v\n", zap.String("bucket", bucket))
		return nil
	}

	var mostRecentObject *s3.Object
	for _, object := range contents {
		if mostRecentObject == nil {
			mostRecentObject = object
		} else {
			if object.LastModified.After(*mostRecentObject.LastModified) {
				mostRecentObject = object
			}
		}
	}

	return Download(bucket, *mostRecentObject.Key)
}

// ListObjects lists the objects in a bucket.
func ListObjects(bucketName string) ([]*s3.Object, error) {
	sess := session.Must(session.NewSession())
	s3cli := s3.New(sess)

	result, err := s3cli.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
	var contents []*s3.Object
	if err != nil {
		log2.Logger().Error("ListObjectsV2", zap.Error(err))
	} else {
		contents = result.Contents
	}
	return contents, err
}
