package main

import (
	"flag"
	"fmt"
	"github.com/hailocab/goamz/s3"
	"github.com/majest/gos3/client"
	"os"
	"strings"
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

	action := os.Args[1]

	switch {
	case action == "cp":
		err = copy(os.Args[2], os.Args[3])
	case action == "ls":

		// check if path is correct. If no arguments are passed
		// list the main directory
		path := ""
		if len(os.Args) > 2 {
			// if argument is passed
			path = os.Args[2]
			// make sure we add slash at the end of the path
			if path[len(path)-1:len(path)] != "/" {
				path = path + "/"
			}
		}

		err = ls(path)

	case action == "rm":
		err = rm(os.Args[2])
	case action == "mv":
		err = mv(os.Args[2], os.Args[3])
	}

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}

func checkFlags() {

	// validate action
	if len(os.Args) <= 1 || contains(actions, os.Args[1]) == -1 {
		fmt.Printf("Invalid action. Available actions are: %+s\n", actions)
		os.Exit(0)
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
		os.Exit(0)
	}
}

func contains(arr []string, search string) int {

	for key, value := range arr {
		if search == value {
			return key
		}
	}

	return -1
}

func rm(fileTarget string) error {
	return s3client.S3Bucket.Del(fileTarget)
}

func copy(fileSource, fileTarget string) error {

	// add required bucket name if bucket name is not provided
	if strings.Index(fileSource, s3client.BucketName) == -1 {
		fileSource = s3client.BucketName + "/" + fileSource
	}

	_, err := s3client.S3Bucket.PutCopy(fileTarget, s3.Private, s3.Options{SSE: true}, fileSource)
	if err != nil {
		return err
	}

	return nil
}

func mv(fileSource, fileTarget string) error {
	err := copy(fileSource, fileTarget)
	if err != nil {
		return err
	}

	return rm(fileSource)
}

func ls(path string) error {

	res, err := s3client.S3Bucket.List(path, "/", "", 1000)

	if err != nil {
		return err
	}

	if res != nil && len(res.CommonPrefixes) > 0 {
		for _, fullpath := range res.CommonPrefixes {
			fmt.Println(fullpath)
		}
	}

	return nil
}
