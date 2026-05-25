// @Author xuanshuiyuan
// HTTP 请求封装包：支持 GET/POST/PostForm/PostJson/PostByte 等请求方式
// 所有请求复用全局 http.Client，自动处理 BOM 头和 JSON 反序列化
package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
)

// httpClient 全局复用的 HTTP 客户端，避免每次请求创建新连接
var httpClient = &http.Client{}

// CurlService HTTP 请求服务，支持链式调用
// 使用方式:
//   result, err := NewCurl().SetUrl("https://api.example.com/data").SetValue(map[string]interface{}{"key": "val"}).PostJson()
type CurlService struct {
	Url     string                 // 请求地址
	Value   map[string]interface{} // 请求参数
	Headers map[string]string      // 自定义请求头
}

// NewCurl 创建 CurlService 实例
func NewCurl() *CurlService {
	return &CurlService{}
}

// SetUrl 设置请求地址（链式调用）
func (c *CurlService) SetUrl(url string) *CurlService {
	c.Url = url
	return c
}

// SetValue 设置请求参数（链式调用）
func (c *CurlService) SetValue(value map[string]interface{}) *CurlService {
	c.Value = value
	return c
}

// SetHeaders 设置自定义请求头（链式调用）
func (c *CurlService) SetHeaders(headers map[string]string) *CurlService {
	c.Headers = headers
	return c
}

func (c *CurlService) Get() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	res, err := httpClient.Get(c.Url)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl get failed, url:%s err:%v", c.Url, err)
		return result, err
	}
	defer res.Body.Close()
	robots, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl get failed, io.ReadAll url:%s, err:%v", c.Url, err)
		return result, err
	}
	robots = bytes.TrimPrefix(robots, []byte("\xef\xbb\xbf"))
	if err = json.Unmarshal(robots, &result); err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl get json unmarshal failed, url:%s, err:%v", c.Url, err)
		return result, err
	}
	return result, nil
}

func (c *CurlService) PostForm() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	DataUrlVal := url.Values{}
	for key, val := range c.Value {
		val1, ok := val.(string)
		if ok {
			DataUrlVal.Add(key, val1)
		}
		val2, ok := val.(float64)
		if ok {
			DataUrlVal.Add(key, goxy.Float64ToString(val2))
		}
	}
	req, err := http.NewRequest("POST", c.Url, strings.NewReader(DataUrlVal.Encode()))
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl post failed, url:%s err:%v", c.Url, err)
		return result, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for key, header := range c.Headers {
		req.Header.Set(key, header)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl post failed, url:%s err:%v", c.Url, err)
		return result, err
	}
	defer resp.Body.Close()
	robots, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl post failed, io.ReadAll url:", c.Url, "err: ", err)
		return result, err
	}
	robots = bytes.TrimPrefix(robots, []byte("\xef\xbb\xbf"))
	if err = json.Unmarshal(robots, &result); err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl post json unmarshal failed, url:%s, err:%v", c.Url, err)
		return result, err
	}
	return result, nil
}

func (c *CurlService) PostJson() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	postString, _ := json.Marshal(c.Value)
	req, err := http.NewRequest("POST", c.Url, strings.NewReader(string(postString)))
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl post failed, url:%s err:%v", c.Url, err)
		return result, err
	}
	req.Header.Add("Content-Type", "application/json")
	for key, header := range c.Headers {
		req.Header.Set(key, header)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl post failed, url:%s err:%v", c.Url, err)
		return result, err
	}
	defer resp.Body.Close()
	robots, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl post failed, io.ReadAll url:", c.Url, "err: ", err)
		return result, err
	}
	robots = bytes.TrimPrefix(robots, []byte("\xef\xbb\xbf"))
	if err = json.Unmarshal(robots, &result); err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl post json unmarshal failed, url:%s, err:%v", c.Url, err)
		return result, err
	}
	return result, nil
}

func (c *CurlService) Post() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	buffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buffer)
	if err := jsonEncoder.Encode(c.Value); err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println(goxy.FmtLog("title.title", "curl post encode failed", "url.title", c.Url, "error.title", err.Error()))
		return result, err
	}
	req, err := http.NewRequest("POST", c.Url, buffer)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println(goxy.FmtLog("title.title", "curl post failed", "url.title", c.Url, "error.title", err.Error()))
		return result, err
	}
	req.Header.Add("Content-Type", "application/json")
	for key, header := range c.Headers {
		req.Header.Set(key, header)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println(goxy.FmtLog("title.title", "curl post", "url.title", c.Url, "error.title", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()
	robots, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl post failed, io.ReadAll url:", c.Url, "err: ", err)
		return result, err
	}
	robots = bytes.TrimPrefix(robots, []byte("\xef\xbb\xbf"))
	if err = json.Unmarshal(robots, &result); err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl post json unmarshal failed, url:%s, err:%v", c.Url, err)
		return result, err
	}
	return result, nil
}

func (c *CurlService) PostByte() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buffer)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(c.Value); err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl post encode failed, url:%s err:%v", c.Url, err)
		return nil, err
	}
	req, err := http.NewRequest("POST", c.Url, buffer)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl post failed, url:%s err:%v", c.Url, err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Printf("curl post failed, url:%s err:%v", c.Url, err)
		return nil, err
	}
	defer resp.Body.Close()
	robots, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl post failed, io.ReadAll url:", c.Url, "err: ", err)
		return nil, err
	}
	return robots, nil
}
