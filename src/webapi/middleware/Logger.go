package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// 2016-09-27 09:38:21.541541811 +0200 CEST
// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700]
// "GET /apache_pb.gif HTTP/1.0" 200 2326
// "http://www.example.com/start.html"
// "Mozilla/4.08 [en] (Win98; I ;Nav)"

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (b bodyLogWriter) Write(bytes []byte) (int, error) {
	b.body.Write(bytes)
	return b.ResponseWriter.Write(bytes)
}

//Logger is the logrus logger handler
func Logger(logger logrus.FieldLogger, notLogged ...string) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, p := range notLogged {
			skip[p] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		requestID := uuid.New()
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID.String())
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		if strings.Contains(path, "/swagger/") {
			return
		}
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()

		if _, ok := skip[path]; ok {
			return
		}

		queryParams := c.Request.URL.Query().Encode()

		originHeaders := map[string]bool{
			"Accept":                    true,
			"Accept-Encoding":           true,
			"Connection":                true,
			"Content-Length":            true,
			"Content-Type":              true,
			"User-Agent":                true,
			"Sec-Fetch-Dest":            true,
			"Accept-Language":           true,
			"Sec-Ch-Ua":                 true,
			"Sec-Ch-Ua-Platform":        true,
			"Sec-Ch-Ua-Mobile":          true,
			"Sec-Fetch-Site":            true,
			"Sec-Fetch-Mode":            true,
			"Sec-Fetch-User":            true,
			"Referer":                   true,
			"Cache-Control":             true,
			"Upgrade-Insecure-Requests": true,
		}

		headers := c.Request.Header
		var headersString string
		var customHeadersString string
		for key, values := range headers {
			for _, value := range values {
				if originHeaders[key] {
					headersString = fmt.Sprint(headersString + key + ":" + value + " ")
				} else {
					customHeadersString = fmt.Sprint(customHeadersString + key + ":" + value + " ")
				}
			}
		}
		domain := c.Request.Host

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		respBody := blw.body.String()

		var respHeadersArr []string
		for header := range c.Writer.Header() {
			respHeadersArr = append(respHeadersArr, header)
		}
		var respHeaders string
		for _, value := range respHeadersArr {
			respHeaders += value + ":" + c.Writer.Header().Get(value) + " "
		}

		entry := logger.WithFields(logrus.Fields{
			"hostname":        hostname,
			"domain":          domain,
			"statusCode":      statusCode,
			"latency":         latency, // time to process
			"clientIP":        clientIP,
			"method":          c.Request.Method,
			"path":            path,
			"referer":         referer,
			"userAgent":       clientUserAgent,
			"queryParams":     queryParams,
			"headers":         headersString,
			"requestBody":     string(bodyBytes),
			"requestID":       requestID,
			"customHeaders":   customHeadersString,
			"responseBody":    respBody,
			"responseHeaders": respHeaders,
		})

		if len(c.Errors) > 0 {
			entry.Warn(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			//msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" %d %d \"%s\" \"%s\" (%dms)", clientIP, hostname, time.Now().Format(timeFormat), c.Request.Method, path, statusCode, dataLength, referer, clientUserAgent, latency)
			if statusCode >= http.StatusInternalServerError {
				entry.Error()
			} else if statusCode >= http.StatusBadRequest {
				entry.Warn()
			} else {
				entry.Info()
			}
		}
	}
}
