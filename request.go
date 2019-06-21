package main

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type logsRequest struct {
	Filter    string `form:"filter"`
	Group     string `form:"group"`
	Stream    string `form:"stream"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Limit     string `form:"limit"`
	NextToken string `form:"next_token"`
}

func logsRequestFromContext(c *gin.Context) (*logsRequest, error) {
	req := &logsRequest{}
	if err := c.Bind(req); err != nil {
		return nil, err
	}
	if req.Group == "" {
		return nil, errors.New("log group must be set")
	}
	return req, nil
}
