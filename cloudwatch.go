package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func newCloudwatchClient(config *Config) (*cloudwatchlogs.CloudWatchLogs, error) {
	shouldConfig := true
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	sess.Config.Region = aws.String(config.Region)

	if config.AccessKey == "" && config.SecretKey == "" && config.Profile == "" {
		log.Println("AWS credentials are not provided, loading from environment...")
		if ec2metadata.New(sess).Available() {
			shouldConfig = false
		} else {
			log.Println("EC2 metadata is not available")
		}
	}

	if config.Profile != "" {
		sess.Config.Credentials = credentials.NewSharedCredentials("", config.Profile)
		shouldConfig = false
	}

	if shouldConfig {
		sess.Config.Credentials = credentials.NewStaticCredentials(
			config.AccessKey,
			config.SecretKey,
			"",
		)
	}

	return cloudwatchlogs.New(sess), nil
}

func (r *logsRequest) cloudwatchInput() *cloudwatchlogs.FilterLogEventsInput {
	input := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: aws.String(r.Group),
	}

	if r.Stream != "" {
		names := []*string{}
		for _, val := range strings.Split(r.Stream, ",") {
			names = append(names, &val)
		}
		input.SetLogStreamNames(names)
	}

	if r.Filter != "" {
		input.SetFilterPattern(r.Filter)
	}

	if r.NextToken != "" {
		input.SetNextToken(r.NextToken)
	}

	if r.StartTime != "" {
		duration, err := time.ParseDuration(r.StartTime)
		if err == nil {
			now := time.Now()
			startTime := now.Add(-duration)
			input.SetStartTime(startTime.Unix() * 1000)
			input.SetEndTime(now.Unix() * 1000)
		} else {
			log.Println("cant parse duration:", err)
		}
	} else {
		now := time.Now()
		input.SetStartTime(now.Add(time.Minute*-15).Unix() * 1000)
		input.SetEndTime(now.Unix() * 1000)
	}

	if r.Limit != "" {
		fmt.Sscanf(r.Limit, "%d", input.Limit)
	} else {
		input.SetLimit(200)
	}

	return input
}
