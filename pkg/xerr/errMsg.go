package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "操作成功"
	message[ServerCommonError] = "服务器开小差了，请稍后再试"
	message[RequestParamError] = "请求参数错误"
	message[TokenExpireError] = "token过期"
	message[TokenGenerateError] = "生成token失败"
	message[DbError] = "数据库繁忙，请稍后再试"
	message[DbUpdateAffectedZeroError] = "更新数据影响行数为0"
	message[DataNoExistError] = "数据不存在"
}

func MapErrMsg(errCode uint32) string {
	if msg, exists := message[errCode]; exists {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}

}

func IsCodeErr(errCode uint32) bool {

	if _, exists := message[errCode]; exists {
		return true
	} else {
		return false
	}
}
