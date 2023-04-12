package proxy

import (
	"log"
	"sync"
	"time"
)

var wg3 sync.WaitGroup
var mux1 sync.Mutex
var ch1 = make(chan int, 50)

var lock2 sync.Mutex
var httpI []ProxyIp
var httpS []ProxyIp

// 初始化
func InitData() {
	//获取配置文件
	GetConfigData()
	//设置线程数量
	ch1 = make(chan int, conf.Config.ThreadNum)
	ch2 = make(chan int, conf.Config.ThreadNum)
	//是否需要抓代理
	if len(ProxyPool) < conf.Config.ProxyNum {
		//抓取代理
		spiderRun()
	}
	//定时判断是否需要获取代理iP
	go func() {
		// 每 60 秒钟时执行一次
		ticker := time.NewTicker(60 * time.Second)
		for range ticker.C {
			if len(ProxyPool) < conf.Config.ProxyNum {
				if !run && !verifyIS {
					log.Printf("代理数量不足 %d\n", conf.Config.ProxyNum)
					//抓取代理
					spiderRun()
				}
			} else {
				//保存代理到本地
				export()
			}
		}
	}()

	// 验证代理存活情况
	go func() {
		verifyTime := time.Duration(conf.Config.VerifyTime)
		ticker := time.NewTicker(verifyTime * time.Second)
		for range ticker.C {
			if !verifyIS && !run {
				VerifyProxy()
			}
		}
	}()
}

func GethttpIp() string {
	lock2.Lock()
	defer lock2.Unlock()
	if len(ProxyPool) == 0 {
		return ""
	}
	for _, v := range ProxyPool {
		if v.Type == "HTTP" {
			is := true
			for _, vv := range httpI {
				if v.Ip == vv.Ip && v.Port == vv.Port {
					is = false
				}
			}
			if is {
				httpI = append(httpI, v)
				return v.Ip + ":" + v.Port
			}
		}
	}
	var addr string
	if len(httpI) != 0 {
		addr = httpI[0].Ip + ":" + httpI[0].Port
	}
	httpI = make([]ProxyIp, 0)

	return addr
}

func GetHttpsIp() string {
	lock2.Lock()
	defer lock2.Unlock()
	if len(ProxyPool) == 0 {
		return ""
	}
	for _, v := range ProxyPool {
		if v.Type == "HTTPS" {
			is := true
			for _, vv := range httpS {
				if v.Ip == vv.Ip && v.Port == vv.Port {
					is = false
				}
			}
			if is {
				httpS = append(httpS, v)
				return v.Ip + ":" + v.Port
			}
		}
	}
	var addr string
	if len(httpS) != 0 {
		addr = httpS[0].Ip + ":" + httpS[0].Port
	}
	httpS = make([]ProxyIp, 0)
	return addr
}
