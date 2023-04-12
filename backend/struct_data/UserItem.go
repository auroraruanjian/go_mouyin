package struct_data

type UserItem struct {
	Uid             string
	Sec_uid         string
	Unique_id       string
	Short_id        string
	Nickname        string
	Signature       string
	Custom_verify   string
	Follower_count  int64 // 粉丝
	Following_count int64 // 关注
	Is_draw         bool  // 是否已抓取
	Aweme_list      map[string]*AwemeItem
}
