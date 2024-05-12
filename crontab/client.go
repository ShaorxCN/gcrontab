package crontab

import (
	"bytes"
	"fmt"
	"gcrontab/constant"
	"gcrontab/utils"
	"io"
	"net/http"
	"time"
)

// ResponseWrapper 返回的包装res
type ResponseWrapper struct {
	StatusCode int
	Body       string
	Header     http.Header
	Start      time.Time
	End        time.Time
}

// Get 获取get 请求
func Get(url string, timeout int, h *http.Header) *ResponseWrapper {
	req, err := http.NewRequest(constant.HTTPMETHODGET, url, nil)
	if err != nil {
		return createRequestError(err)
	}

	if h != nil {
		req.Header = *h
	}

	return request(req, timeout)
}

// PostForm 发送post Form 请求
func PostForm(url, params string, timeout int, h *http.Header) *ResponseWrapper {
	buf := bytes.NewBufferString(params)
	req, err := http.NewRequest(constant.HTTPMETHODPOST, url, buf)
	if err != nil {
		return createRequestError(err)
	}

	if h != nil {
		req.Header = *h
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	return request(req, timeout)
}

// PostJSON 发送post请求
func PostJSON(url, body string, timeout int, h *http.Header) *ResponseWrapper {
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return createRequestError(err)
	}
	if h != nil {
		req.Header = *h
	}
	req.Header.Set("Content-type", "application/json")

	return request(req, timeout)
}

func request(req *http.Request, timeout int) *ResponseWrapper {
	wrapper := &ResponseWrapper{StatusCode: -1, Body: "", Header: make(http.Header)}
	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Millisecond
	}
	start := utils.Now()
	wrapper.Start = start
	resp, err := client.Do(req)
	if err != nil {
		wrapper.Body = fmt.Sprintf("执行HTTP请求错误-%s", err.Error())
		wrapper.End = utils.Now()
		return wrapper
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		wrapper.Body = fmt.Sprintf("读取HTTP请求返回值失败-%s", err.Error())
		wrapper.End = utils.Now()
		return wrapper
	}
	end := utils.Now()
	wrapper.StatusCode = resp.StatusCode
	wrapper.Body = string(body)
	wrapper.End = end
	return wrapper
}

func createRequestError(err error) *ResponseWrapper {
	errorMessage := fmt.Sprintf("创建HTTP请求错误-%s", err.Error())
	end := utils.Now()
	return &ResponseWrapper{StatusCode: -1, Body: errorMessage, End: end}
}
