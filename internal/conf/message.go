// @Author xuanshuiyuan
package conf

var MessagePushType = map[string]string{
	"SmsLoginCode": "AlismsGgc-LoginCode",
}

var WxApp = map[string]map[string]interface{}{
	"Ggc": map[string]interface{}{
		"AppId":              "wxfc6eb5bbcde5f947",
		"Secret":             "4e89b9ce408a4246d2af549e24d7f3c0",
		"MiniprogramAppid":   "wx1aacb2cf610ca17e",
		"GetAccessTokenCall": "GgcWxAppAccessToken",
		"VehicleEvaluationOrderAuditResults": map[string]interface{}{
			"Code":                "T2-xaCpcVphRZjjDVlx9SpYm__AU1BwuWudEwChgcNs",
			"MiniprogramPagepath": "/pages/content/index",
			"Message":             "车辆估价报告生成通知", //
		},
	},
}

var Sms = map[string]map[string]interface{}{

}
