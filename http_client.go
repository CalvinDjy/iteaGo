package itea

import (
	"bytes"
	"context"
	"github.com/CalvinDjy/iteaGo/ilog"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HttpClient struct {
	Ctx context.Context
	debug bool
}

func (c *HttpClient) Construct() {
	c.debug = c.Ctx.Value(DEBUG).(bool)
}

func (c *HttpClient) Get(u string, h map[string]string, host string) (result []byte, err error) {
	if c.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【GET请求】耗时：", time.Since(start), ", 请求地址[", u,"]")
		}()
	}

	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}

	//client := &http.Client{Transport: tr}
	client := &http.Client{}

	req, err := http.NewRequest("GET", u, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(host, "") {
		req.Host = host
	}

	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *HttpClient) Post(u string, p map[string]string, h map[string]string, host string) (result []byte, err error) {
	if c.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【POST请求】耗时：", time.Since(start), ", 请求地址[", u, "]")
		}()
	}

	postParams := url.Values{}
	for k, v := range p {
		postParams.Set(k, v)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", u, strings.NewReader(postParams.Encode()))
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(host, "") {
		req.Host = host
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *HttpClient) PostFile(u string, file string, filekey string, p map[string]string, h map[string]string, host string) (result []byte, err error) {
	if c.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【POST FILE请求】耗时：", time.Since(start), ", 请求地址[", u, "]")
		}()
	}

	//创建一个缓冲区对象,后面的要上传的body都存在这个缓冲区里
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile(filekey, filepath.Base(file))
	if err != nil {
		return nil, err
	}

	//打开文件
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	//把文件流写入到缓冲区里去
	_, err = io.Copy(fileWriter, f)
	if err != nil {
		return nil, err
	}

	for k, v := range p {
		bodyWriter.WriteField(k, v)
	}

	contentType := bodyWriter.FormDataContentType()

	bodyWriter.Close()

	client := &http.Client{}

	req, err := http.NewRequest("POST", u, ioutil.NopCloser(bodyBuf))
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(host, "") {
		req.Host = host
	}

	req.Header.Set("Content-Type", contentType)

	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}