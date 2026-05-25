package production

import (
	logic "go_project_template/internal"
)

func NewConfig() *logic.ConfigConf {
	config := &logic.ConfigConf{
		Env:          Env,
		Verification: Verification,
		ApiTimer:     ApiTimer,
	}
	return config
}

const Env = "production"

// 01 web, 10 wxapp
// 01 B, 10 C
var Verification = &logic.Verification{
	SourceList:           []string{"0101", "1010", "0110"},
	SourceRedisList:      map[string]string{"0101": "admin_info", "1010": "userinfo"},
	SourceExplainList:    map[string]string{"0101": "网页B端", "1010": "微信小程序C端", "0110": "网页C端"},
	SourceEngExplainList: map[string]string{"0101": "Web-B", "1010": "Wxapp-C", "0110": "Web-C"},
	KeyList:              map[string]string{"0101": "PODV1034UWMYOPAMZQ1G", "1010": "JQ6L16WSSN4L9E08FO5S", "0110": "V5FS5AZLW68RXNUMK69R"},
	EditionList:          map[string]int64{"v1.0.0": 100},
}

var ApiTimer = &logic.ApiTimer{
	Ip: []string{"127.0.0.1"},
}
