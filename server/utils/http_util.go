package utils

import (
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

var httpClient *http.Client

func init() {
	// 失败可重试客户端
	retryClient := retryablehttp.NewClient()
	//重试不超过50次
	retryClient.RetryMax = 1
	retryClient.RetryWaitMin = 10 * time.Millisecond
	retryClient.RetryWaitMax = 50 * time.Millisecond
	httpClient = retryClient.StandardClient()
	httpClient.Timeout = time.Minute * 5
	httpClient.Transport = &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 50,
		DisableKeepAlives:   false,
	}
}

func Get(url string, params map[string]string) (statusCode int, body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}
	// construct query params
	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	return doSend(req)
}

func NormalGetStr(url string) (resp string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println(err)
		return "", err
	}

	if err != nil {
		return "", err
	}

	_, body, err := doSend(req)
	if err != nil {
		return "", err
	}

	// 将响应体转换为字符串
	bodyString := string(body)
	return bodyString, err
}

func GetStr(url string) (resp string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println(err)
		return "", err
	}
	req.Header.Set("Connection", "keep-alive") //设置请求头
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	if err != nil {
		return "", err
	}

	_, body, err := doSend(req)
	if err != nil {
		return "", err
	}

	// 将响应体转换为字符串
	bodyString := string(body)
	return bodyString, err
}

func GetWithCookie(url string, cookies []*http.Cookie) (resp string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println(err)
		return "", err
	}
	req.Header.Set("Connection", "keep-alive") //设置请求头
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	for _, c := range cookies {
		req.AddCookie(c)
	}

	if err != nil {
		return "", err
	}

	_, body, err := doSend(req)
	if err != nil {
		return "", err
	}

	// 将响应体转换为字符串
	bodyString := string(body)
	return bodyString, err
}

func GetWithHeaders(url string, params map[string]string, headers map[string]string) (statusCode int, body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}
	// construct query params
	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	// add headers
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return doSend(req)
}

func Post(url string, sJSON string) (statusCode int, body []byte, err error) {
	// create new http request
	req := &http.Request{}
	if sJSON != "" {
		req, err = http.NewRequest("POST", url, strings.NewReader(sJSON))
	} else {
		req, err = http.NewRequest("POST", url, nil)
	}

	if err != nil {
		return 0, nil, err
	}
	// add headers
	req.Header.Add("Content-type", "application/json; charset=utf-8")
	return doSend(req)
}

func PostWithHeaders(url string, sJSON string, headers map[string]string) (statusCode int, body []byte, err error) {
	// create new http request
	req, err := http.NewRequest("POST", url, strings.NewReader(sJSON))
	if err != nil {
		return 0, nil, err
	}
	// add headers
	req.Header.Add("Content-type", "application/json; charset=utf-8")

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return doSend(req)
}

func PostWithCookie(url string, sJSON string, c *http.Cookie) (statusCode int, body []byte, err error) {
	// create new http request
	req, err := http.NewRequest("POST", url, strings.NewReader(sJSON))
	if err != nil {
		return 0, nil, err
	}
	// add headers
	req.AddCookie(c)
	req.Header.Add("Content-type", "application/json; charset=utf-8")
	return doSend(req)
}

func doSend(req *http.Request) (statusCode int, body []byte, err error) {
	// do
	res, err := httpClient.Do(req)
	if err != nil {
		switch tmpErr := err.(type) {
		case net.Error:
			if tmpErr.Timeout() {
				return 0, nil, err
			}
		}
		return 0, nil, err
	}
	// read response body
	bodyBytes, err := io.ReadAll(res.Body)
	// close the connection to reuse it
	defer func() { _ = res.Body.Close() }()
	if err != nil {
		return 0, bodyBytes, err
	}
	if res.StatusCode != http.StatusOK {
		return res.StatusCode, bodyBytes, err
	}
	return res.StatusCode, bodyBytes, nil
}

func ParseUrlParams(url string) map[string]string {
	// 解析URL参数
	params := make(map[string]string)
	if strings.Contains(url, "?") {
		queryParams := strings.Split(url, "?")[1]
		for _, param := range strings.Split(queryParams, "&") {
			kv := strings.Split(param, "=")
			params[kv[0]] = kv[1]
		}
	}
	return params
}
