package itea

import (
	"net/http"
	"io/ioutil"
	"strings"
	"net/url"
	"time"
	"github.com/CalvinDjy/iteaGo/ilog"
)

type HttpClient struct {

}

func (c *HttpClient) Get(u string, h map[string]string, host string) (result []byte, err error) {
	start := time.Now()
	defer func() {
		ilog.Info("【GET请求】耗时：", time.Since(start), ", 请求地址[", u,"]")
	}()

	client := &http.Client{}

	req, err := http.NewRequest("GET", u, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(host, "") {
		req.Host = host
	}

	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *HttpClient) Post(u string, p map[string]string, h map[string]string, host string) (result []byte, err error) {
	start := time.Now()
	defer func() {
		ilog.Info("【POST请求】耗时：", time.Since(start), ", 请求地址[", u,"]")
	}()
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

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
