package middleware

import (
	"gcrontab/utils"
	"math"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		unescapeURL, _ := url.PathUnescape(c.Request.URL.String())
		path := c.Request.URL.Path
		start := utils.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / float64(time.Millisecond)))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknow"
		}
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		entry := logrus.NewEntry(log).WithFields(logrus.Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"url":        unescapeURL,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
			"startTime":  start,
		})
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			format := "[%s] \"%s %s\" %d %d (%dms)"
			args := []interface{}{start, c.Request.Method, path, statusCode, dataLength, latency}
			switch {
			case statusCode > 499:
				entry.Errorf(format, args...)
			case statusCode > 399:
				entry.Warnf(format, args...)
			default:
				entry.Infof(format, args...)
			}
		}
	}
}
