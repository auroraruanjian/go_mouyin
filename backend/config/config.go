package config

import (
	"net/http"
)

// 接口KEY
const APIKEY = ""

// 是否测试模式
const ISTEST = false

// 最大抓取的视频,发布环境不可设置太大生成的Excel过大
var MAX_VIDEO_COUNT = 200

// 启动的抓取的线程数量
const WORKER_COUNT = 20

var (
	// RequestSleepTime 重试间隔时间 如1000为1秒
	RequestSleepTime = 1000
	// RequestAttempts 重试次数
	RequestAttempts    = 5
	RequestUserAgent   = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36`
	RequestRediract    = true
	RequestGlobeCookie = false
	RequestUseProxy    = false
)

var CookieJar http.CookieJar

/*
type ProxyIp struct {
	Ip         string //IP地址
	Port       string //代理端口
	Country    string //代理国家
	Province   string //代理省份
	City       string //代理城市
	Isp        string //IP提供商
	Type       string //代理类型
	Anonymity  string //代理匿名度, 透明：显示真实IP, 普匿：显示假的IP, 高匿：无代理IP特征
	Time       string //代理验证
	Speed      string //代理响应速度
	SuccessNum int    //验证请求成功的次数
	RequestNum int    //验证请求的次数
	Source     string //代理源
}

var ProxyPool []ProxyIp

func init() {
	fmt.Println("开始导入代理IP")
	//导入代理缓存
	file, err := os.OpenFile("./backend/config/data.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("代理json文件打开错误：" + err.Error())
		return
	}
	defer file.Close()
	all, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("代理json解析错误：" + err.Error())
		return
	}
	if len(all) == 0 {
		fmt.Println("代理json为空")
		return
	}
	err = json.Unmarshal(all, &ProxyPool)
	if err != nil {
		fmt.Println("代理json解析错误：" + err.Error())
		return
	}
}
*/
