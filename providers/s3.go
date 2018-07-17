package providers

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3 uploads the resulting artifacts to S3
type S3 struct{}

func upload(ctx ProviderContext, file string) string {
	f, err := os.Open(file)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	conf := aws.Config{
		Region: aws.String(ctx.Config["aws_region"]),
	}

	s := session.New(&conf)

	svc := s3manager.NewUploader(s)

	location := ctx.Config["location"]

	if location == "" {
		location = "evacuate.tar.gz"
	} else {
		location = strings.Replace(location, "%epoch%", strconv.FormatInt(time.Now().Unix(), 10), -1)

		hostname, err := os.Hostname()

		if err != nil {
			panic(err)
		}

		location = strings.Replace(location, "%hostname%", hostname, -1)

		if !strings.HasSuffix(location, ".tar.gz") {
			location += ".tar.gz"
		}
	}

	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(ctx.Config["aws_bucket"]),
		Key:    aws.String(location),
		Body:   bufio.NewReader(f),
	})

	if err != nil {
		ctx.Logger.Errorf("Error upload to S3: %s", err)
		os.Exit(1)
	}

	return result.Location
}

// Run TODO
func (p S3) Run(ctx ProviderContext, file string) {
	ctx.Finish <- upload(ctx, file)
}
