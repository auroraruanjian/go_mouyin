package util

import (
	"crypto/tls"
	"douyin/backend/config"
	"douyin/backend/proxy"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func SendRetryGet(request_url string, header map[string]string) ([]byte, error) {
	return SendRequestGetBody("GET", request_url, "", header)
}

func SendRetryPost(url string, reqBody string, header map[string]string) ([]byte, error) {
	return SendRequestGetBody("POST", url, reqBody, header)
}

func SendRequestGetBody(method string, postUrl string, reqBody string, header map[string]string) ([]byte, error) {
	// 发送Http请求
	_, body, err := SendRetryRequest(method, postUrl, reqBody, header)
	if err != nil {
		return nil, err
	}

	return body, nil
}

/**
* 可进行重试的http-post请求方法
* method GET,POST
* url 发送请求的路由
* reqBody 请求参数
 */
func SendRetryRequest(method string, url string, reqBody string, header map[string]string) (*http.Response, []byte, error) {
	errString := ""
	for index := 0; index < config.RequestAttempts; index++ {

		res, body, err := SendRequest(method, url, reqBody, header)
		if err == nil {
			return res, body, nil
		}
		if index != 0 {
			errString += "|\n" + err.Error()
		} else {
			errString += err.Error()
		}

		if body != nil {
			// 如果返回body，接口有可能限流，休眠重试
			sleepTime := time.Duration(2*index+1) * time.Millisecond * time.Duration(config.RequestSleepTime)
			time.Sleep(sleepTime)

			//time.Millisecond*1000 * time.Duration(2*3+1)
			fmt.Printf("休眠时间 %v \n", sleepTime)
		} else {
			time.Sleep(time.Millisecond * time.Duration(config.RequestSleepTime))
		}
	}

	return nil, nil, errors.New("SendRetry err:" + errString)
}

func SendRequest(method string, postUrl string, reqBody string, header map[string]string) (*http.Response, []byte, error) {
	netTransport := &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, time.Second*5) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			_ = conn.SetDeadline(time.Now().Add(time.Second * 15)) //设置发送接受数据超时
			return conn, nil
		},
		ResponseHeaderTimeout: time.Second * 15,
		MaxIdleConnsPerHost:   10,
	}

	// 如果启动代理
	if config.RequestUseProxy {
		len_proxypool := len(proxy.ProxyPool)
		if len_proxypool > 0 {
			proxy, _ := url.Parse("http://" + proxy.GethttpIp())
			netTransport.Proxy = http.ProxyURL(proxy)
			netTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			fmt.Printf("代理服务开启：%v\n", proxy.String())
		} else {
			fmt.Println("代理池为空")
		}
	}

	//可以通过client中transport的Dail函数,在自定义Dail函数里面设置建立连接超时时长和发送接受数据超时
	client := &http.Client{
		Transport: netTransport,
	}

	if config.RequestGlobeCookie {
		if config.CookieJar == nil {
			jarObject, err := cookiejar.New(&cookiejar.Options{})
			if err != nil {
				log.Printf("Globe Cookie: %v", err)
				return nil, nil, errors.New("Globe Cookie err:" + err.Error())
			}
			config.CookieJar = jarObject
		}

		client.Jar = config.CookieJar
	}

	var requestDo *http.Request
	var err error
	if reqBody != "" {
		requestDo, err = http.NewRequest(method, postUrl, strings.NewReader(reqBody))
	} else {
		requestDo, err = http.NewRequest(method, postUrl, nil)
	}

	if err != nil {
		log.Printf("NewRequest error: %v", err)
		return nil, nil, errors.New("httpPost err:" + err.Error())
	}

	//提交请求;用指定的方法，网址，可选的主体放回一个新的*Request
	for key, value := range header {
		if key == "user-agent" {
			continue
		}
		requestDo.Header.Add(key, value)
	}
	requestDo.Header.Add("user-agent", config.RequestUserAgent)
	//requestDo.Header.Add("Content-Type", "application/json")

	// 是否启动301，302跳转
	if !config.RequestRediract {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	//前面预处理一些参数，状态，Do执行发送；处理返回结果;Do:发送请求,
	res, err := client.Do(requestDo)
	if nil != err {
		log.Printf("http[%v] error: %v, url: %v, params: %v\n", method, err, postUrl, reqBody)
		return nil, nil, errors.New("http err:" + err.Error())
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if nil != err {
		log.Println("ReadAll err:", err)
		return res, nil, errors.New("ReadAll err:" + err.Error())
	}

	fmt.Printf("请求URL:%#v\n", postUrl)
	fmt.Printf("请求头:%#v\n", res.Request.Header)
	fmt.Println("响应状态:" + res.Status)
	fmt.Printf("响应头:%#v\n", res.Header)
	fmt.Printf("响应包长度:%#v\n", len(data))
	//fmt.Printf("响应数据:%#v\n", string(data))

	if len(data) == 0 {
		log.Println("Empty Data Error:", postUrl)
		return res, data, errors.New("api return empty data")
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusFound ||
			res.StatusCode == http.StatusMovedPermanently {
			log.Println("Response Status Code:" + res.Status)
			return res, data, nil
		} else {
			log.Println("Response Status Code:" + res.Status)
			return nil, nil, errors.New("Response Status Code:" + res.Status)
		}
	}

	return res, data, nil
}
