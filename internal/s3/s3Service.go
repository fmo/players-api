package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/fmo/players-api/config"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Service struct {
	Session s3iface.S3API
}

func NewS3Service() (*Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.GetAwsRegion())},
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		Session: s3.New(sess),
	}, nil
}

func (s Service) Save(s3Key, url string) error {
	s3Bucket := config.GetS3Bucket()

	if s.checkImageAlreadyUploaded(s3Bucket, s3Key) == true {
		return nil
	}

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Upload the file to S3
	object, err := s.Session.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"object": object.String(),
	}).Debug("Photo uploaded")

	return nil
}

func (s Service) checkImageAlreadyUploaded(s3Bucket, s3Key string) bool {
	objectMetadata, err := s.Session.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Key),
	})

	if err == nil {
		log.WithFields(log.Fields{
			"s3Key":    s3Key,
			"metadata": objectMetadata.Metadata,
		}).Debugf("Image is already upload for %v", s3Key)

		return true
	}

	return false
}
