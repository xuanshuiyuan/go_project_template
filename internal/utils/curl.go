// @Author xuanshuiyuan
package utils

import (
	"bytes"
	"encoding/json"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type CurlService struct {
	Url     string
	Value   map[string]interface{}
	Headers map[string]string
}

//var Curl *CurlService

func NewCurl() *CurlService {
	return &CurlService{}
}

func (c *CurlService) SetUrl(url string) *CurlService {
	c.Url = url
	return c
}

func (c *CurlService) SetValue(value map[string]interface{}) *CurlService {
	c.Value = value
	return c
}

func (c *CurlService) SetHeaders(headers map[string]string) *CurlService {
	c.Headers = headers
	return c
}

// @Title Get
// @Description Get请求
// @Author xuanshuiyuan 2021/12/29 10:48:00
// @Param
// @Return map[string]interface{}, error
func (c *CurlService) Get() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	res, err := http.Get(c.Url)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl get failed, url:%s err:%v", c.Url, err)
		return result, err
	}
	defer res.Body.Close()
	robots, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl get failed, ioutil.ReadAll url:%s, err:%v", c.Url, err)
		return result, err
	}
	robots = bytes.TrimPrefix(robots, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}
	json.Unmarshal(robots, &result)
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
	res, err := http.NewRequest("POST", c.Url, strings.NewReader(DataUrlVal.Encode()))
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl post failed, url:%s err:%v", c.Url, err)
		return result, err
	}
	res.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if len(c.Headers) > 0 {
		for key, header := range c.Headers {
			res.Header.Set(key, header)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
	}
	defer resp.Body.Close()
	robots, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl post failed, ioutil.ReadAll url:", c.Url, "err: ", err)
		return result, err
	}
	robots = bytes.TrimPrefix(robots, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}
	json.Unmarshal(robots, &result)
	return result, nil
}

func (c *CurlService) PostJson() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	postString, _ := json.Marshal(c.Value)
	res, err := http.NewRequest("POST", c.Url, strings.NewReader(string(postString)))
	//res, err := http.NewRequest("POST", c.Url, postString)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl post failed, url:%s err:%v", c.Url, err)
		return result, err
	}
	res.Header.Add("Content-Type", "application/json")
	if len(c.Headers) > 0 {
		for key, header := range c.Headers {
			res.Header.Set(key, header)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
	}
	defer resp.Body.Close()
	robots, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl post failed, ioutil.ReadAll url:", c.Url, "err: ", err)
		return result, err
	}
	robots = bytes.TrimPrefix(robots, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}
	json.Unmarshal(robots, &result)
	return result, nil
}

// @Title Post
// @Description Post请求
// @Author xuanshuiyuan 2021/12/29 10:48:00
// @Param
// @Return map[string]interface{}, error
func (c *CurlService) Post() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	buffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buffer)
	//jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(c.Value)
	postString := buffer
	res, err := http.NewRequest("POST", c.Url, postString)
	//log.Info(conf.Config.Base.LogFileName, "").Println(goxy.FmtLog("url.title", c.Url, "parmas.title", c.Value))
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println(goxy.FmtLog("title.title", "curl post failed", "url.title", c.Url, "error.title", err.Error()))
		return result, err
	}
	res.Header.Add("Content-Type", "application/json")
	if len(c.Headers) > 0 {
		for key, header := range c.Headers {
			res.Header.Set(key, header)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println(goxy.FmtLog("title.title", "curl post", "url.title", c.Url, "error.title", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()
	robots, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl post failed, ioutil.ReadAll url:", c.Url, "err: ", err)
		return result, err
	}
	robots = bytes.TrimPrefix(robots, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}
	json.Unmarshal(robots, &result)
	//log.Info(conf.Config.Base.LogFileName, "").Println(goxy.FmtLog("title.title", "curl post", "url.title", c.Url, "params.title", c.Value, "res.title", result))
	return result, nil
}

// @Title Post
// @Description Post请求
// @Author xuanshuiyuan 2021/12/29 10:48:00
// @Param
// @Return map[string]interface{}, error
func (c *CurlService) PostByte() ([]byte, error) {
	var result = []byte{}
	buffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buffer)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(c.Value)
	postString := buffer
	res, err := http.NewRequest("POST", c.Url, postString)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl post failed, url:%s err:%v", c.Url, err)
		return result, err
	}
	res.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
	}
	defer resp.Body.Close()
	robots, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("curl post failed, ioutil.ReadAll url:", c.Url, "err: ", err)
		return result, err
	}
	return robots, nil
}
