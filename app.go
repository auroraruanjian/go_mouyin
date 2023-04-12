package main

import (
	"context"
	"douyin/backend/api"
	"douyin/backend/config"
	"douyin/backend/struct_data"
	"douyin/backend/util"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	api *api.Douyin
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.api = api.NewDouyinApi(ctx)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) DrawData(inputUrl string) struct_data.UserItem {
	runtime.LogInfo(a.ctx, "DrawData is run")

	sec_uid := "MS4wLjABAAAAP4sRCear3dNc1yvDmhDBgwCUem8WFWxbwosy0LSXREfuRrSr5qJ65LPVA1d4BB_X"

	// API
	user_item, err := a.api.GetDetail(sec_uid)
	if err != nil {
		fmt.Print(err)
		return user_item
	}
	fmt.Printf("%#v \n", user_item)
	//time.Sleep(time.Second * 2)

	max_cursor, aweme_item, err := a.api.GetAweme(sec_uid, "0")
	if err != nil {
		fmt.Print(err)
		return user_item
	}
	fmt.Printf("max_cursor:%#v \n aweme_item:%#v \n", max_cursor, aweme_item)

	/*
		follow_map := &sync.Map{}
		max_time, err := a.api.GetFollow("d5a7966cc2c766e5a0a50f63cf7f4cc9", follow_map, sec_uid, "0")
		if err != nil {
			fmt.Print(err)
			return user_item
		}
		fmt.Printf("max_time:%#v \n follow_map:%#v \n", max_time, follow_map)
		follow_map.Range(func(key any, value any) bool {
			runtime.LogDebugf(a.ctx, "%#v \n %#v \n", key, value)
			return true
		})
	*/

	runtime.LogInfo(a.ctx, "DrawData is finsh")
	return user_item
}

func (a *App) SearchUser(keyword string) map[string]string {
	// error_message
	user_item, err := a.api.SearchUser(keyword)
	err_message := ""
	if err != nil {
		err_message = err.Error()
	}

	json_string, error := json.Marshal(user_item)
	if error != nil {
		err_message = error.Error()
	}
	return map[string]string{
		"user_item": string(json_string),
		"error":     err_message,
	}
}

func (a *App) GetFollow(sessionid string, sec_uid string, max_time string) map[string]string {
	// error_message
	res_max_time, user_item, err := a.api.GetFollow(sessionid, sec_uid, max_time)
	err_message := ""
	if err != nil {
		err_message = err.Error()
	}

	json_string, error := json.Marshal(user_item)
	if error != nil {
		err_message = error.Error()
	}
	return map[string]string{
		"max_time":  res_max_time,
		"user_item": string(json_string),
		"error":     err_message,
	}
}

func (a *App) GetAweme(sec_uid string, max_cursor string) map[string]string {
	// error_message
	max_cursor, ameme_list, err := a.api.GetAweme(sec_uid, max_cursor)
	err_message := ""
	if err != nil {
		err_message = err.Error()
	}

	json_string, error := json.Marshal(ameme_list)
	if error != nil {
		err_message = error.Error()
	}
	return map[string]string{
		"max_cursor": max_cursor,
		"ameme_list": string(json_string),
		"error":      err_message,
	}
}

func (a *App) GetCode() map[string]string {
	// error_message
	qr_json, err := a.api.GetQrcode()
	err_message := ""
	if err != nil {
		err_message = err.Error()
	}

	return map[string]string{
		"json":  qr_json,
		"error": err_message,
	}
}

func (a *App) CheckQrcode(token string) map[string]string {
	// error_message
	json, err := a.api.CheckQrcode(token)
	err_message := ""
	if err != nil {
		err_message = err.Error()
	}
	fmt.Println(json, err)

	return map[string]string{
		"json":  json,
		"error": err_message,
	}
}

func (a *App) LogRedirect(url string) map[string]string {
	// error_message
	json, sessionid, err := a.api.LogRedirect(url)
	err_message := ""
	status := "0"
	if err != nil {
		err_message = err.Error()
	} else {
		status = "1"
		if sessionid != "" {
			query_json, error1 := a.api.GetQueryUser(sessionid)
			if error1 != nil {
				err_message = error1.Error()
			} else {
				json = query_json
				status = "2"
			}
		}
	}
	fmt.Println(json, err)

	return map[string]string{
		"status":    status, // 0:获取session id失败，1:获取sessionid成功，但获取用户信息失败，2:获取用户信息成功
		"json":      json,
		"sessionid": sessionid,
		"error":     err_message,
	}
}

// 获取用户ID
func (a *App) GetQueryUser(sessionid string) map[string]string {
	// error_message
	json, err := a.api.GetQueryUser(sessionid)
	err_message := ""
	if err != nil {
		err_message = err.Error()
	}
	fmt.Println(json, err)

	return map[string]string{
		"json":  json,
		"error": err_message,
	}
}

func (a *App) GetRelation(sessionid string) map[string]string {
	// error_message
	json, err := a.api.GetRelation(sessionid)
	err_message := ""
	if err != nil {
		err_message = err.Error()
	}
	fmt.Println(json, err)

	return map[string]string{
		"json":  json,
		"error": err_message,
	}
}

/*
抓取数据
*/
func (a *App) DrawData1(inputUrl string) (sync.Map, bool) {
	var user_list = sync.Map{}

	// "https://www.douyin.com/user/MS4wLjABAAAAP4sRCear3dNc1yvDmhDBgwCUem8WFWxbwosy0LSXREfuRrSr5qJ65LPVA1d4BB_X?showTab=post"
	// 解析sec_user_id
	re := regexp.MustCompile(`/user/([\w].*?)\?`)
	restr := re.FindStringSubmatch(inputUrl)

	if len(restr) < 2 {
		util.WriteError("error", "未匹配到数据")
		return user_list, false
	}

	sec_uid := restr[1]
	util.WriteError("tip", "获取到sec_uid:"+sec_uid)

	util.WriteError("tip", "开始抓取用户数据：:"+sec_uid)

	self_user, err := a.api.GetDetail(sec_uid)
	if err != nil {
		util.WriteError("error", fmt.Sprintf("%#v \n", err))
		return user_list, false
	}
	util.WriteError("tip", "抓取到用户数据："+self_user.Nickname)

	util.WriteError("tip", "开始抓取用户关注列表")
	min_time := ""
	i := 1
	for {
		resposne_min_time, _, err := a.api.GetFollow("d5a7966cc2c766e5a0a50f63cf7f4cc9", self_user.Uid, min_time)
		if err != nil {
			util.WriteError("error", err.Error())
			break
		}

		util.WriteError("tip", fmt.Sprintf("第 %d 页抓取完成 \n", i))
		if resposne_min_time == "0" {
			util.WriteError("error", "无下一页数据，停止抓取")
			break
		}

		min_time = resposne_min_time
		time.Sleep(time.Second * 1)
		i++
	}

	util.WriteError("tip", "用户关注列表抓取完成")

	if config.ISTEST {
		config.MAX_VIDEO_COUNT = 20
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////////
	// 任务计数器
	taskWg := sync.WaitGroup{}

	// 定义静态变量 用于外部访问内部方法
	// 启动线程池 len=channel通道容量，超过容量生产者阻塞，容量变成0 消费者阻塞
	pool := util.NewWrokPool(config.WORKER_COUNT, 5).Start()
	//////////////////////////////////////////////////////////////////////////////////////////////////////

	user_list.Range(func(uid, item interface{}) bool {
		user_floow_item := item.(*struct_data.UserItem)

		taskWg.Add(1)
		task := &util.Task{
			Args: uid,
			F: func(args ...interface{}) {
				defer taskWg.Done() // 任务完成计数器减一

				//////////////////////////////////////////////////////////////////////////////////////////////////////
				util.WriteError("tip", fmt.Sprintf("开始抓取用户【%v】视频列表 \n", user_floow_item.Nickname))
				start_cursor := "0"
				start_page_number := 1
				for {
					util.WriteError("info", fmt.Sprintf("开始抓取用户【%v】-【%v】\n", user_floow_item.Nickname, start_cursor))
					max_cursor, aweme_list, _ := a.api.GetAweme(user_floow_item.Sec_uid, start_cursor)

					if aweme_list == nil {
						util.WriteError("info", "数据为空，停止抓取")
						break
					}

					for _, range_AwemeItem := range aweme_list {
						user_floow_item.Aweme_list[range_AwemeItem.Aweme_id] = &range_AwemeItem
					}

					util.WriteError("info", fmt.Sprintf("成功抓取第 %d 页数据\n", start_page_number))

					if max_cursor == "0" || strings.TrimSpace(max_cursor) == "" {
						user_floow_item.Is_draw = true
						util.WriteError("info", "无下一页数据，停止抓取")
						break
					}

					if config.ISTEST {
						user_floow_item.Is_draw = true
						util.WriteError("info", "测试环境跳过用户，停止抓取")
						break
					}

					if len(user_floow_item.Aweme_list) >= config.MAX_VIDEO_COUNT {
						user_floow_item.Is_draw = true
						util.WriteError("info", "超过最大视频抓取数量，跳过")
						break
					}

					if start_page_number > config.MAX_VIDEO_COUNT/8 {
						user_floow_item.Is_draw = true
						util.WriteError("info", "异常抓取，跳过")
						break
					}

					start_cursor = max_cursor
					time.Sleep(time.Second * 1)
					start_page_number++
				}
				util.WriteError("tip", fmt.Sprintf("完成抓取用户【%v】视频列表 \n", user_floow_item.Nickname))
				//////////////////////////////////////////////////////////////////////////////////////////////////////
			},
		}
		pool.PushTask(task)
		return true
	})

	util.WriteError("tip", "任务入队完成")
	//等待任务执行完成
	taskWg.Wait()

	// 回收资源
	close(pool.Pool)
	pool.Stop()
	util.WriteError("tip", "任务全部执行完成")

	return user_list, true
}
