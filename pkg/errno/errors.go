package errno

type EmptyErr struct {
	msg string
}

func NewEmptyErr() EmptyErr {
	return EmptyErr{
		msg: "not find data",
	}
}
func NewParamErr() EmptyErr {
	return EmptyErr{
		msg: "param error",
	}
}
func (e EmptyErr) Error() string {
	return e.msg
}

func (e EmptyErr) RuntimeError() {
}
