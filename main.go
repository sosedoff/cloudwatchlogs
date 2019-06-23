package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
	flags "github.com/jessevdk/go-flags"
)

func serveLogStream(c *gin.Context) {
	client := clientFromContext(c)

	req, err := logsRequestFromContext(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	output, err := client.FilterLogEvents(req.cloudwatchInput())
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, output)
}

func serverGroupsList(c *gin.Context) {
	client := clientFromContext(c)

	output, err := client.DescribeLogGroups(nil)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, output.LogGroups)
}

func serveStreamsList(c *gin.Context) {
	client := clientFromContext(c)

	groupName := c.Query("group")
	if groupName == "" {
		c.JSON(400, gin.H{"error": "group name is not provided"})
		return
	}

	output, err := client.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(groupName),
		OrderBy:      aws.String("LastEventTime"),
	})
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, output.LogStreams)
}

func clientFromContext(c *gin.Context) *cloudwatchlogs.CloudWatchLogs {
	return c.MustGet("cloudwatch").(*cloudwatchlogs.CloudWatchLogs)
}

func serveHome(c *gin.Context) {
	client := clientFromContext(c)

	// Fetch available log groups
	output, err := client.DescribeLogGroups(nil)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.HTML(200, "/static/index.html", gin.H{
		"log_groups": output.LogGroups,
	})
}

func serveStaticAsset(c *gin.Context) {
	asset, ok := Assets.Files[c.Request.URL.Path]
	if !ok {
		c.AbortWithStatus(404)
		return
	}
	c.Data(200, "text/plain", asset.Data)
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".html") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func setGinDefaults() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	log.SetFlags(log.LstdFlags)
}

func main() {
	config, err := readConfig()
	if err != nil {
		switch err.(type) {
		case *flags.Error:
			// no need to print error, flags package already does that
		default:
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	client, err := newCloudwatchClient(config)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := client.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{}); err != nil {
		log.Fatal(err)
	}

	tpl, err := loadTemplate()
	if err != nil {
		log.Fatal(err)
	}

	setGinDefaults()

	router := gin.Default()
	router.SetHTMLTemplate(tpl)

	if config.AuthUser != "" && config.AuthPassword != "" {
		router.Use(gin.BasicAuth(gin.Accounts{
			config.AuthUser: config.AuthPassword,
		}))
	}

	router.Use(func(c *gin.Context) {
		c.Set("cloudwatch", client)
	})

	router.GET("/", serveHome)
	router.GET("/static/:file", serveStaticAsset)
	router.GET("/groups", serverGroupsList)
	router.GET("/streams", serveStreamsList)
	router.POST("/logs", serveLogStream)

	go func() {
		exec.Command("open", "http://"+config.ListenAddr()).Run()
	}()

	if err := router.Run(config.ListenAddr()); err != nil {
		log.Fatal(err)
	}
}
