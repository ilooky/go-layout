package errno

var (
	OK        = NewResp(1, "OK")
	ErrServer = NewResp(0, "服务异常，请联系管理员")
	ErrParam  = NewResp(0, "参数有误")
)
