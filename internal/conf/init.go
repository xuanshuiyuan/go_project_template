// @Author  xuanshuiyuan
package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	xtime "github.com/Terry-Mao/goim/pkg/time"
	"github.com/golang/freetype/truetype"
	logic "go_project_template/internal"
	"go_project_template/internal/conf/develop"
	"go_project_template/internal/conf/local"
	"go_project_template/internal/conf/production"
	"net"
	"os"
)

//ConfigService 定义了基本配置
type ConfigService struct {
	Base    *BaseInfo
	Iris    *IrisInfo
	Doctron *Doctron
	Redis   *Redis
	Mysql   *Mysql
	Mongodb *Mongodb
	Conf    *logic.ConfigConf
	Font    *truetype.Font
}

//IrisInfo iris基本配置
type IrisInfo struct {
	DisablePathCorrection             bool
	EnablePathEscape                  bool
	FireMethodNotAllowed              bool
	DisableBodyConsumptionOnUnmarshal bool
	Charset                           string
	Port                              int32
	HtmlTemplate                      string
}

type Doctron struct {
	Hostname string
	Username string
	Password string
	Port     int32
}

//Mysql mysql数据库的配置
type Mysql struct {
	Hostname  string
	Username  string
	Password  string
	Port      int32
	DataBase  string
	Drive     string
	Network   string
	Charset   string
	TimeZone  string
	DbLogMode bool
}

//Redis redis的配置
type Redis struct {
	Network      string
	Addr         string
	Auth         string
	Active       int
	Idle         int
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	IdleTimeout  xtime.Duration
	Expire       xtime.Duration
	DataBase     int
}

//Mongodb mongodb的配置
type Mongodb struct {
	Hostname   string
	Username   string
	Password   string
	Port       int32
	ConnectStr string
}

type Nsq struct {
	NsqdAddress string
	TopicName   string
	ChannelName string
}

//BaseInfo 基本配置
type BaseInfo struct {
	LogPath     string
	LogFileName string
	TtfPath     string
}

var ConfPath string
var Config *ConfigService

func ConfInit() {
	flag.StringVar(&ConfPath, "conf", "", "")
	flag.Parse()
	Init()
}

// @Title Init
// @Description 配置文件初始化
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param
// @Return err
func Init() (err error) {
	Config = DefaultConf()
	_, err = toml.DecodeFile(ConfPath, &Config)
	return
}

func DefaultConf() *ConfigService {
	return &ConfigService{
		Conf: LoadConfigConfig(),
	}
}

// @Title LoadConfigConfig
// @Description 根据环境变量加载配置
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param
// @Return core.ConfigConf
func LoadConfigConfig() *logic.ConfigConf {
	var res = &logic.ConfigConf{}
	env := os.Getenv("ENV")
	if env == "" {
		panic("缺少环境变量")
	}
	ip, err := GetLocalIP()
	if err != nil {
		panic(err)
	}
	switch env {
	case "local":
		r := local.NewConfig()
		r.LocalIP = ip
		return r
	case "develop":
		r := develop.NewConfig()
		r.LocalIP = ip
		return r
	case "production":
		r := production.NewConfig()
		r.LocalIP = ip
		return r
	}
	return res
}

func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}

//初始化字体
//func InitFont() {
//	fontBytes, err := ioutil.ReadFile(Config.Base.TtfPath)
//	if err != nil {
//		panic("载入字体失败" + err.Error())
//	}
//	//载入字体数据
//	Config.Font, err = freetype.ParseFont(fontBytes)
//	if err != nil {
//		panic("载入字体失败" + err.Error())
//	}
//}
