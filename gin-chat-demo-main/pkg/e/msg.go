package e

var MsgFlags = map[int]string{
	SUCCESS:               "ok",
	UpdatePasswordSuccess: "修改密码成功",
	NotExistInentifier:    "该第三方账号未绑定",
	ERROR:                 "fail",
	InvalidParams:         "请求参数错误",
	ErrorDatabase:         "数据库操作出错,请重试",

	ErrorHeaderData: "缺少Authorization数据",
	ErrorValidToken: "无效的Token",
	ErrorGetFile:    "从网络获取文件失败",
	ErrorCreateFile: "创建文件目录失败",
	ErrorPlaceFile:  "文件放置失败",
	ErrorCreateData: "创建数据失败",
	ErrorGroupFull:  "群人数已满",
	ErrorNoData:     "无法获取群数据",
	ExistUser:       "用户已在群中",

	WebsocketSuccessMessage: "解析content内容信息",
	WebsocketSuccess:        "发送信息，请求历史纪录操作成功",
	WebsocketEnd:            "请求历史纪录，但没有更多记录了",
	WebsocketOnlineReply:    "针对回复信息在线应答成功",
	WebsocketOfflineReply:   "针对回复信息离线回答成功",
	WebsocketLimit:          "请求收到限制",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
