// @Author xuanshuiyuan
package models

type MessageQueue struct {
	Id           int64  `json:"id"`
	Mobile       string `json:"mobile"`
	Type         int8   `json:"type"`
	Title        string `json:"title"`
	TemplateCode string `json:"template_code"`
	Params       string `json:"params"`
	Status       int8   `json:"status"`
	Reason       string `json:"reason"`
	CreateTime   int64  `json:"create_time"`
}

type MessagePushQueue struct {
	Id           int64  `json:"id"`
	Sender       string `json:"sender"`
	Receiver     string `json:"receiver"`
	Channel      int8   `json:"channel"`
	Type         string `json:"type"`
	TemplateCode string `json:"template_code"`
	Content      string `json:"content"`
	Params       string `json:"params"`
	Status       int8   `json:"status"`
	Reason       string `json:"reason"`
	CreateTime   int64  `json:"create_time"`
}
