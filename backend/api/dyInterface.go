package api

import (
	"douyin/backend/struct_data"
	"sync"
)

type DyInterface interface {
	// 获取用户详情信息
	// @return user_item,error
	GetDetail(sec_uid string) (struct_data.UserItem, error)

	// 获取用户关注用户-授权接口，需要sessionid
	// @return min_time,error
	GetFollow(user_list *sync.Map, user_id string, max_time string) (string, error)

	// 获取用户发布视频列表
	// @return max_cursor,aweme_item,error
	GetAweme(sec_uid string, max_cursor string) (string, []struct_data.AwemeItem, error)

	GetQrcode() (string, error)

	CheckQrcode(token string) (string, error)

	LogRedirect(rediract_url string) (string, string, error)

	SearchUser(keyword string) ([]struct_data.UserItem, error)
}

var DouYinAPI = map[string]string{
	"detail":          "https://www.douyin.com/aweme/v1/web/user/profile/other/",
	"follow":          "https://www.douyin.com/aweme/v1/web/user/following/list/",
	"aweme":           "https://www.douyin.com/aweme/v1/web/aweme/post/",
	"search":          "https://www.douyin.com/aweme/v1/web/discover/search/",
	"get_qrcode":      "https://sso.douyin.com/get_qrcode/",
	"check_qrconnect": "https://sso.douyin.com/check_qrconnect/",
	"query_user":      "https://www.douyin.com/aweme/v1/web/query/user/",
	"relation":        "https://www.douyin.com/aweme/v1/web/im/spotlight/relation/",
}

var WhoSeCardAPI = map[string]string{
	"detail": "https://api.whosecard.com/api/douyin/aweme/user/detail",
	"follow": "https://api.whosecard.com/api/douyin/aweme/user/following/list",
	"aweme":  "https://api.whosecard.com/api/douyin/aweme/post",
}
