package xerr

var CodeTest = map[int]string{
	SERVER_COMMON_ERROR: "服务异常，请稍后再试",
	REQUEST_PARAM_ERROR: "参数不正确",
	TOKEN_EXPIRE_ERROR:  "token已过期，请重新登录",
	DB_ERROR:            "数据库繁忙，请稍后再试",
}

func ErrMsg(errcode int) string {
	if msg, ok := CodeTest[errcode]; ok {
		return msg
	}
	return CodeTest[SERVER_COMMON_ERROR]
}
