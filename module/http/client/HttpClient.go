package httpClient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpParams struct {
	// 用于设置header头
	Header http.Header
	// 用于设置url参数
	Params url.Values
	// 用于设置form参数
	Forms url.Values
}

type HttpClient struct {
	// 定义的http客户端
	Client http.Client
}

// GET get请求
func (httpClient *HttpClient) GET(url string, params HttpParams) (string, error) {
	// 1. 设置header头
	return httpClient._send("GET", url, params)
}

func (httpClient *HttpClient) GETResp(url string, params HttpParams) (*http.Response, error) {
	return httpClient._sendResp("GET", url, params)
}

func (httpClient *HttpClient) POSTResp(url string, params HttpParams) (*http.Response, error) {
	return httpClient._sendResp("POST", url, params)
}

// POST post请求
func (httpClient *HttpClient) POST(url string, params HttpParams) (string, error) {
	return httpClient._send("POST", url, params)

}

func (httpClient *HttpClient) _sendResp(method string, url string, params HttpParams) (*http.Response, error) {

	// 构造request
	req, _ := http.NewRequest(method, url, strings.NewReader(params.Forms.Encode()))
	fmt.Println(req)
	// 检查query参数是否为空 不为空的话增加query参数
	if len(params.Params) > 0 {
		req.URL.RawQuery = params.Params.Encode()
	}
	// 设置headers
	req.Header = params.Header
	// 检查header是否为空
	rep, err := httpClient.Client.Do(req)
	if err == nil {

		// 如果状态码 大于等于400 请求就算是失败 error不为空
		if rep.StatusCode >= 400 {
			return rep, errors.New("http status is error")
		}

		// 正常请求 返回body
		return rep, err
	} else {

		return nil, err
	}

}

// 最底层的请求
func (httpClient *HttpClient) _send(method string, url string, params HttpParams) (string, error) {

	// 构造request
	req, _ := http.NewRequest(method, url, strings.NewReader(params.Forms.Encode()))
	fmt.Println(req)
	// 检查query参数是否为空 不为空的话增加query参数
	if len(params.Params) > 0 {
		req.URL.RawQuery = params.Params.Encode()
	}
	// 设置headers
	req.Header = params.Header
	// 检查header是否为空
	rep, err := httpClient.Client.Do(req)
	if err == nil {

		// 如果状态码 大于等于400 请求就算是失败 error不为空
		body, bodyError := ioutil.ReadAll(rep.Body)
		// 判断body是否为空 不为空返回
		if bodyError != nil {
			return "read body error", errors.New("http parse body error")
		}

		if rep.StatusCode >= 400 {
			fmt.Println(string(body))
			return string(body), errors.New("http status is error")
		}

		// 正常请求 返回body
		return string(body), err
	} else {

		return "http error", err
	}

}

func (httpClient HttpClient) GetDefaultClient() *HttpClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	return &HttpClient{
		Client: client,
	}
}
