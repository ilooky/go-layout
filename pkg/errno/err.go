package errno

import (
	"github.com/ilooky/go-layout/pkg/guava/json"
	"github.com/ilooky/logger"
)

var _ Resp = (*err)(nil)

type Resp interface {
	// i 为了避免被其他包实现
	i()
	WithData(data interface{}) Resp
	WithStacks(err error) Resp
	WithMessage(message string) Resp
	WithID(id string) Resp
	ToString() string
}

type err struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
	Stacks  string      `json:"stacks,omitempty"`
	ID      string      `json:"id,omitempty"` // 当前请求的唯一ID，便于问题定位，忽略也可以
}

func NewResp(code int, msg string) Resp {
	return &err{
		Code:    code,
		Message: msg,
		Content: nil,
	}
}

func Ok() Resp {
	return NewResp(1, "")
}

func ParamErr() Resp {
	return NewResp(0, "参数有误")
}

func ServerErr() Resp {
	return NewResp(0, "服务异常，请联系管理员")
}

func (e *err) i() {}

func (e *err) GetData() interface{} {
	return e.Content
}

func (e *err) WithData(data interface{}) Resp {
	e.Content = data
	return e
}

func (e *err) WithMessage(message string) Resp {
	e.Message = message
	return e
}

func (e *err) WithStacks(err error) Resp {
	logger.Error(err)
	e.Stacks = err.Error()
	return e
}

func (e *err) WithID(id string) Resp {
	e.ID = id
	return e
}

func (e *err) ToString() string {
	err := &struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
		ID   string      `json:"id,omitempty"`
	}{
		Code: e.Code,
		Msg:  e.Message,
		Data: e.Content,
		ID:   e.ID,
	}

	raw, _ := json.Marshal(err)
	return string(raw)
}
