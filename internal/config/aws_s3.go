package config

import (
	"fmt"
	"strings"
)

type AwsS3 struct {
	Region          string `mapstructure:"region" yaml:"region"`
	AccessKeyId     string `mapstructure:"access-key-id" yaml:"access-key-id"`
	AccessKeySecret string `mapstructure:"access-key-secret" yaml:"access-key-secret"`
	Bucket          string `mapstructure:"bucket" yaml:"bucket"`
	PreviewUrl      string `mapstructure:"preview-url" yaml:"preview-url"`
}

func (c AwsS3) FormatUrl(objectKey string) string {
	if objectKey == "" {
		return ""
	}
	if strings.Index(objectKey, "http") < 0 {
		//return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.Bucket, c.Region, objectKey)
		return fmt.Sprintf("%s/%s", c.PreviewUrl, objectKey)
	}
	return objectKey
}
