package api

import (
	"context"
	"douyin/backend/config"
	"douyin/backend/struct_data"
	"douyin/backend/util"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Douyin struct {
	ctx        context.Context
	ttwid      string
	ttwid_time time.Time
}

// NewApp creates a new App application struct
func NewDouyinApi(ctx context.Context) *Douyin {
	return &Douyin{
		ctx: ctx,
	}
}

func (d Douyin) api(baseUrl string, params map[string]string, header map[string]string) ([]byte, error) {
	parse_url, _ := url.Parse(baseUrl)

	// 构建参数
	url_object := url.Values{}
	for key, value := range params {
		url_object.Set(key, value)
	}
	parse_url.RawQuery = url_object.Encode()

	request_url := parse_url.String()

	if baseUrl != DouYinAPI["get_qrcode"] &&
		baseUrl != DouYinAPI["check_qrconnect"] &&
		params != nil {
		// 生成xbogus
		xbogus, execJsError := util.ExecJsXbogus(parse_url.RawQuery, config.RequestUserAgent)
		if execJsError != nil {
			runtime.LogError(d.ctx, execJsError.Error())
			return nil, errors.New("JS解析错误")
		}
		request_url = request_url + "&X-Bogus=" + xbogus
	}

	if baseUrl == DouYinAPI["check_qrconnect"] {
		config.RequestGlobeCookie = true
	} else {
		config.RequestGlobeCookie = false
	}

	fmt.Printf("%v %v \n", baseUrl, DouYinAPI["aweme"])
	if baseUrl == DouYinAPI["aweme"] {
		config.RequestUseProxy = true
	} else {
		config.RequestUseProxy = false
	}

	if header == nil {
		header = map[string]string{}
	}
	//value, ok := map[key]
	if _, ok := header["referer"]; !ok {
		header["referer"] = "https://www.douyin.com/"
	}

	//runtime.LogDebug(d.ctx, request_url+"&X-Bogus="+xbogus)
	//"https://www.douyin.com/aweme/v1/web/user/profile/other/?device_platform=webapp&aid=6383&sec_user_id=MS4wLjABAAAAP4sRCear3dNc1yvDmhDBgwCUem8WFWxbwosy0LSXREfuRrSr5qJ65LPVA1d4BB_X&X-Bogus=DFSzswSLGwtANGyRtc3p-z9WcBno"
	response_byte, err := util.SendRetryGet(request_url, header)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return nil, errors.New("接口请求异常：" + request_url)
	}
	//status_code
	status_code := gjson.GetBytes(response_byte, "status_code").Int()
	if status_code != 0 {
		runtime.LogDebug(d.ctx, string(response_byte))
		return nil, errors.New("返回数据错误：" + string(response_byte))
	}

	return response_byte, nil
}

// 获取用户详情信息
// @return user_item,error
func (d Douyin) GetDetail(sec_uid string) (struct_data.UserItem, error) {
	var userItem struct_data.UserItem

	ttwid, err := d.GetTtwid(true)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return userItem, err
	}

	request_header := map[string]string{
		"cookie": fmt.Sprintf("ttwid=%v; douyin.com; ", ttwid),
	}

	//fmt.Println("DEB: | ttwid:" + ttwid)
	//fmt.Println("DEB: | cookie:" + request_header["cookie"])

	request_params := map[string]string{
		"aid":             "6383",
		"device_platform": "webapp",
		"sec_user_id":     sec_uid,
	}

	response_byte, err := d.api(DouYinAPI["detail"], request_params, request_header)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return userItem, err
	}
	//runtime.LogDebug(d.ctx, string(response_byte))

	user := gjson.GetBytes(response_byte, "user")

	userItem = struct_data.UserItem{
		Uid:             user.Get("uid").String(),
		Sec_uid:         user.Get("sec_uid").String(),
		Unique_id:       user.Get("unique_id").String(),
		Short_id:        user.Get("short_id").String(),
		Nickname:        user.Get("nickname").String(),
		Signature:       user.Get("signature").String(),
		Custom_verify:   user.Get("custom_verify").String(),
		Follower_count:  user.Get("follower_count").Int(),
		Following_count: user.Get("following_count").Int(),
		Is_draw:         false,
		Aweme_list:      make(map[string]*struct_data.AwemeItem),
	}

	return userItem, nil
}

// 获取用户关注用户-授权接口，需要sessionid
// @return min_time,error
func (d Douyin) GetFollow(sessionid string, sec_uid string, max_time string) (string, []struct_data.UserItem, error) {
	user_item := []struct_data.UserItem{}

	request_params := map[string]string{
		"aid":         "6383",
		"sec_user_id": sec_uid,
		"max_time":    max_time,
		"msToken":     util.DouyinMsToken(107),
	}

	request_header := map[string]string{
		"cookie": "sessionid=" + sessionid,
	}

	response_byte, err := d.api(DouYinAPI["follow"], request_params, request_header)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return "", user_item, err
	}

	gjson.GetBytes(response_byte, "followings").ForEach(func(key, value gjson.Result) bool {
		user_item = append(user_item, struct_data.UserItem{
			Uid:             value.Get("uid").String(),
			Sec_uid:         value.Get("sec_uid").String(),
			Unique_id:       value.Get("unique_id").String(),
			Short_id:        value.Get("short_id").String(),
			Nickname:        value.Get("nickname").String(),
			Signature:       value.Get("signature").String(),
			Custom_verify:   value.Get("custom_verify").String(),
			Follower_count:  value.Get("follower_count").Int(),
			Following_count: value.Get("following_count").Int(),
			Is_draw:         false,
			Aweme_list:      make(map[string]*struct_data.AwemeItem),
		})
		return true
	})

	response_min_time := gjson.GetBytes(response_byte, "result.min_time").String()

	return response_min_time, user_item, nil
}

// 获取用户发布视频列表
// @return max_cursor,aweme_item,error
func (d Douyin) GetAweme(sec_uid string, max_cursor string) (string, []struct_data.AwemeItem, error) {
	var aweme_list []struct_data.AwemeItem

	request_params := map[string]string{
		"aid":         "6383",
		"sec_user_id": sec_uid,
		"max_cursor":  max_cursor,
		"count":       "100",
		//"msToken":     util.DouyinMsToken(107),
	}

	request_header := map[string]string{
		"referer": "https://www.douyin.com/user/" + sec_uid,
	}

	response_byte, err := d.api(DouYinAPI["aweme"], request_params, request_header)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return "", aweme_list, err
	}
	runtime.LogDebug(d.ctx, string(response_byte))
	gjson.GetBytes(response_byte, "aweme_list").ForEach(func(key, value gjson.Result) bool {
		aweme_list = append(aweme_list, struct_data.AwemeItem{
			Aweme_id:       value.Get("aweme_id").String(),
			Uid:            value.Get("author.uid").String(),
			Desc:           value.Get("desc").String(),
			Address_info:   value.Get("poi_info.address_info").String(),
			Share_url:      value.Get("share_url").String(),
			Ip_attribution: value.Get("ip_attribution").String(),
		})
		return true
	})

	response_max_cursor := gjson.GetBytes(response_byte, "max_cursor").String()

	return response_max_cursor, aweme_list, nil
}

func (d *Douyin) GetTtwid(fource bool) (string, error) {
	// 如果没有wwid或过期
	if d.ttwid == "" || !d.ttwid_time.After(time.Now()) || fource {
		request_url := "https://ttwid.bytedance.com/ttwid/union/register/"

		request_params := `
			{
				"region": "cn",
				"aid": 1768,
				"needFid": false,
				"service": "www.ixigua.com",
				"migrate_info": { "ticket": "", "source": "node" },
				"cbUrlProtocol": "https",
				"union": true
			}
		`
		res, _, err := util.SendRetryRequest("POST", request_url, request_params, nil)

		if err != nil {
			runtime.LogError(d.ctx, err.Error())
			return "", err
		}

		for _, c := range res.Cookies() {
			if c.Name == "ttwid" {
				d.ttwid = c.Value
				d.ttwid_time = time.Now().Add(time.Minute)

				return c.Value, nil
			}
		}

		return "", errors.New("请求错误 未获取ttwid")

	}

	return d.ttwid, nil
}

// 获取用户发布视频列表
// @return body,error
func (d Douyin) GetQrcode() (string, error) {
	response_byte, err := d.api(DouYinAPI["get_qrcode"], nil, nil)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return "", err
	}

	return string(response_byte), nil
}

// 获取用户发布视频列表
// @return body,sessionid,error
func (d Douyin) CheckQrcode(token string) (string, error) {
	ttwid, err := d.GetTtwid(false)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return "", err
	}

	request_header := map[string]string{
		"cookie": fmt.Sprintf("ttwid=%v; douyin.com; ", ttwid),
	}

	request_params := map[string]string{
		"token": token,
	}

	response_byte, err := d.api(DouYinAPI["check_qrconnect"], request_params, request_header)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return "", err
	}

	return string(response_byte), nil
}

// 获取用户发布视频列表
// @return body,sessionid,error
func (d Douyin) LogRedirect(rediract_url string) (string, string, error) {
	_, err := d.GetTtwid(false)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return "", "", err
	}

	request_header := map[string]string{
		//"Cookie": fmt.Sprintf("ttwid=%v; douyin.com; ", ttwid),
	}

	res, body, err := util.SendRetryRequest("GET", rediract_url, "", request_header)

	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return "", "", err
	}

	sessionid := ""
	for _, c := range res.Cookies() {
		if c.Name == "sessionid" {
			sessionid = c.Value
		}
	}
	// 如果
	if sessionid == "" {
		for _, c := range res.Request.Cookies() {
			if c.Name == "sessionid" {
				sessionid = c.Value
			}
		}
	}

	return string(body), sessionid, nil
}

// 获取用户发布视频列表
// @return body,sessionid,error
func (d Douyin) SearchUser(keyword string) ([]struct_data.UserItem, error) {
	user_item := []struct_data.UserItem{}
	request_params := map[string]string{
		"aid":     "6383",
		"keyword": keyword,
		"count":   "10",
	}

	request_header := map[string]string{
		"cookie": fmt.Sprintf("msToken=%v; ", util.DouyinMsToken(107)),
	}

	response_byte, err := d.api(DouYinAPI["search"], request_params, request_header)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return user_item, err
	}

	gjson.GetBytes(response_byte, "user_list").ForEach(func(key, value gjson.Result) bool {
		user_item = append(user_item, struct_data.UserItem{
			Uid:             value.Get("user_info.uid").String(),
			Sec_uid:         value.Get("user_info.sec_uid").String(),
			Unique_id:       value.Get("user_info.unique_id").String(),
			Short_id:        value.Get("user_info.short_id").String(),
			Nickname:        value.Get("user_info.nickname").String(),
			Signature:       value.Get("user_info.signature").String(),
			Custom_verify:   value.Get("user_info.custom_verify").String(),
			Follower_count:  value.Get("user_info.follower_count").Int(),
			Following_count: value.Get("user_info.following_count").Int(),
			Is_draw:         false,
			Aweme_list:      make(map[string]*struct_data.AwemeItem),
		})
		return true
	})

	return user_item, nil
}

// 获取用户基本信息
// @return body,error
func (d Douyin) GetQueryUser(sessionid string) (string, error) {
	ttwid, err := d.GetTtwid(false)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return "", err
	}

	request_header := map[string]string{
		"cookie": fmt.Sprintf("ttwid=%v; douyin.com; sessionid=%v;", ttwid, sessionid),
	}

	request_params := map[string]string{
		"device_platform": "webapp",
	}

	response_byte, err := d.api(DouYinAPI["query_user"], request_params, request_header)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return "", err
	}

	return string(response_byte), nil
}

// 获取用户关系
// @return body,error
func (d Douyin) GetRelation(sessionid string) (string, error) {
	ttwid, err := d.GetTtwid(false)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return "", err
	}

	request_header := map[string]string{
		"cookie": fmt.Sprintf("ttwid=%v; douyin.com; sessionid=%v;", ttwid, sessionid),
	}

	//request_params := map[string]string{}

	response_byte, err := d.api(DouYinAPI["relation"], nil, request_header)
	if err != nil {
		runtime.LogError(d.ctx, err.Error())
		return "", err
	}

	return string(response_byte), nil
}
