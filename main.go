package main

import (
	"flag"
	"fmt"
	"github.com/hailocab/goamz/s3"
	"github.com/majest/gos3/client"
	"os"
)

var bucket = flag.String("bucket", "", "S3 bucket to use")
var awsKey = flag.String("awsKey", "", "AWS Key")
var awsSecret = flag.String("awsSecret", "", "AWS Secret")
var awsRegionName = flag.String("awsRegion", "eu-west-1", "AWS region")
var help = flag.Bool("help", false, "Prints this help")

var s3client *client.S3Client

// available actions
var actions = []string{"ls", "cp", "mv", "rm"}

func main() {

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	checkFlags()

	var err error
	s3client, err = client.NewClient(*bucket, *awsKey, *awsSecret, *awsRegionName)

	// s3 error
	if err != nil {
		fmt.Printf("%s", err.Error())
		os.Exit(1)
	}

	switch {
	case os.Args[1] == "cp":
		copy(os.Args[2], os.Args[3])
	}

}

func checkFlags() {

	// validate action
	if contains(actions, os.Args[1]) == -1 {
		fmt.Printf("Invalid action. Available actions are: %+s\n", actions)
		os.Exit(1)
	}

	if *awsKey == "" {
		*awsKey = os.Getenv("AWS_ACCESS_KEY")
	}

	if *awsSecret == "" {
		*awsSecret = os.Getenv("AWS_SECRET_KEY")
	}

	if *awsSecret == "" || *awsKey == "" {
		fmt.Println("Missing AWS Key or Secret")
		os.Exit(1)
	}

	if *bucket == "" {
		*bucket = os.Getenv("AWS_BUCKET")
	}

	if *bucket == "" {
		fmt.Printf("Missing bucket name: %s", *bucket)

		os.Exit(1)
	}

	fmt.Println(len(os.Args), os.Args)

}

func contains(arr []string, search string) int {

	for key, value := range arr {
		if search == value {
			return key
		}
	}

	return -1
}

func copy(fileTarget, fileSource string) {
	_, err := s3client.S3Bucket.PutCopy(fileTarget, s3.Private, s3.Options{SSE: true}, fileSource)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
