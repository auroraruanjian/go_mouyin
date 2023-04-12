package util

import (
	"douyin/backend/struct_data"
	"fmt"
	"os/user"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

func Write2File(user_list sync.Map) {
	WriteError("tip", "表格生成中...")
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			WriteError("error", err.Error())
		}
	}()
	streamWriter, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		WriteError("error", err.Error())
		return
	}
	styleID, err := file.NewStyle(&excelize.Style{Font: &excelize.Font{Color: "777777"}})
	if err != nil {
		WriteError("error", err.Error())
		return
	}

	// 设置宽度
	if err := streamWriter.SetColWidth(1, 2, 12); err != nil {
		WriteError("error", err.Error())
		return
	}
	if err := streamWriter.SetColWidth(5, 5, 50); err != nil {
		WriteError("error", err.Error())
		return
	}
	if err := streamWriter.SetColWidth(7, 7, 100); err != nil {
		WriteError( "error", err.Error())
		return
	}

	if err := streamWriter.SetRow("A1",
		[]interface{}{
			excelize.Cell{StyleID: styleID, Value: "昵称"},
			excelize.Cell{StyleID: styleID, Value: "抖音号"},
			excelize.Cell{StyleID: styleID, Value: "认证信息"},
			excelize.Cell{StyleID: styleID, Value: "粉丝数"},
			excelize.Cell{StyleID: styleID, Value: "个性签名"},
			excelize.Cell{StyleID: styleID, Value: "IP属地"},
			excelize.Cell{StyleID: styleID, Value: "视频文字"},
			excelize.Cell{StyleID: styleID, Value: "短ID"},
			excelize.Cell{StyleID: styleID, Value: "用户ID"},
			excelize.Cell{StyleID: styleID, Value: "加密ID"},
		},
		excelize.RowOpts{Height: 30, Hidden: false}); err != nil {
		WriteError("error", err.Error())
		return
	}

	//WriteError(text_grid, "tip", fmt.Sprintf("抓取到数据 %#v\n", user_list))

	rowID := 2
	user_list.Range(func(uid, item interface{}) bool {
		user_item := item.(*struct_data.UserItem)
		WriteError("info", fmt.Sprintf("循环数据 %#v\n", user_item.Nickname))

		excel_desc := ""
		excel_ip := ""

		for _, aweme_item := range user_item.Aweme_list {
			excel_desc += aweme_item.Desc
			if aweme_item.Ip_attribution != "" {
				excel_ip = aweme_item.Ip_attribution
			}
		}

		username := user_item.Unique_id
		if username == "" {
			username = user_item.Short_id
		}

		row := make([]interface{}, 50)
		row[0] = user_item.Nickname
		row[1] = username
		row[2] = user_item.Custom_verify
		row[3] = user_item.Follower_count
		row[4] = user_item.Signature
		row[5] = excel_ip
		row[6] = excel_desc
		row[7] = user_item.Short_id
		row[8] = user_item.Uid
		row[9] = user_item.Sec_uid

		cell, _ := excelize.CoordinatesToCellName(1, rowID)
		if err := streamWriter.SetRow(cell, row, excelize.RowOpts{Height: 20, Hidden: false}); err != nil {
			WriteError("error", err.Error())
			return false
		}
		rowID++
		return true
	})

	if err := streamWriter.Flush(); err != nil {
		WriteError( "error", err.Error())
		return
	}

	c_user, err := user.Current()
	if err != nil {
		fmt.Print(err)
	}

	time_now := time.Now().Format("20060102-150405")
	if err := file.SaveAs(c_user.HomeDir + "/Desktop/抓取数据" + time_now + ".xlsx"); err != nil {
		WriteError("error", err.Error())
	}

	WriteError("tip", "表格生成完成")
}


