package machine

/*
错误类型
 */


type ErrType int
const (
	TMPLERR ErrType = iota
	SENDERR
)

type Err struct {
	ErrMsg string
	ErrType ErrType
	error
}
func (e Err)String() string{
	return e.ErrMsg
}
func (e Err)Error() string{
	return e.String()
}

type TmplErr struct {
	Err
}
func NewTmplErr(err string) TmplErr{
	e := TmplErr{}
	e.ErrType = TMPLERR
	e.ErrMsg = err
	return e
}

type SendErr struct {
	Err
}
func NewSendErr(err string) SendErr {
	e := SendErr{}
	e.ErrMsg = err
	e.ErrType = SENDERR
	return e
}

