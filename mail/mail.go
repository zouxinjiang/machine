package mail

import (
	"machine"
	"net/smtp"
	"strings"
	"encoding/json"
	)
type tmpl string
const (
	textTmpl tmpl = "To:<${TO}$>\r\n"+
"From:<${FROM}$>\r\n"+
"Subject:${SUBJECT}$\r\n\r\n"+
"${CONTENT}$"
)
type MsgType int
const (
	TEXT MsgType = iota
)

func (mt MsgType)ToString() string{
	return mt.String()
}

func (mt MsgType)String() string{
	switch mt {
	case TEXT:
		return "TEXT"
	default:
		return "unkown msg type"
	}
}


type Machine struct {
	auth smtp.Auth
	msg string
	msgType machine.MsgType
	username string
	password string
	host string
	prefixVar string
	tmpl map[string]string
	port string
	to []string

}

func NewMachine(url string) Machine{
	m := Machine{}
	m.port = "25"
	m.host = url
	m.prefixVar = "$"
	m.tmpl = make(map[string]string)
	m.to = make([]string,0)
	return m
}

func (m Machine) SendMsg() (int, error) {
	l := len(m.msg)
	err := smtp.SendMail(m.host+":"+m.port,m.auth,m.username,m.to,[]byte(m.msg))
	return l,err
}

func (m *Machine) RegisteTmpl(name, tmpl string) error {
	m.tmpl[name] = tmpl
	return nil
}

func (m *Machine) BuildTxtMsg(tmplName string, data map[string]string) error {
	res,ok := m.tmpl[tmplName]
	if !ok {
		return machine.NewTmplErr("模板"+tmplName+"不存在")
	}
	for k,v := range data {
		res = strings.Replace(res,m.prefixVar+k,v,-1)
	}
	m.msg = res
	return nil
}

func (Machine) BuildStreamMsg(tmplName string, data interface{}) error {
	panic("implement me")
}

func (m *Machine) SetPrefixVar(prefix string) {
	m.prefixVar = prefix
}

func (m *Machine) GenerateTmpl(msgType machine.MsgType, data map[string]interface{}) (tpl string, curerr error) {
	defer func() {
		if e := recover(); e != nil {
			tpl = ""
			emsg := e.(error).Error()
			if strings.Contains(emsg,"interface {} is nil, not string") {
				curerr = machine.NewTmplErr("请传入所需参数")
			}else{
				curerr = machine.NewTmplErr("其他错误。"+"参数类型不正确")
			}
		}
	}()

	//将data的键全部变大写
	for k,v := range data{
		data[strings.ToUpper(k)] = v
		if k != strings.ToUpper(k){
			delete(data,k)
		}
	}

	m.msgType = msgType
	switch msgType {
	case TEXT:
		tpl = strings.Replace(string(textTmpl),"${FROM}$",m.username,-1)
		tpl = strings.Replace(tpl,"${SUBJECT}$",data["SUBJECT"].(string),-1)
		tpl = strings.Replace(tpl,"${CONTENT}$",data["CONTENT"].(string),-1)
		tt,ok  := data["TO"].([]string)
		m.to = tt
		if !ok {
			curerr = machine.NewTmplErr("类型错误。data['TO']类型应该为[]string")
		}
		js,err := json.Marshal(tt)
		str := string(js)
		curerr = err
		tpl = strings.Replace(tpl,"${TO}$",str[1:len(str)-1],-1)
	default:
		return "",machine.NewTmplErr("未知模板类型")
	}

	return tpl,curerr
}

func (m * Machine) Login(data map[string]string) error {
	m.username = data["username"]
	m.password = data["password"]
	m.auth = smtp.PlainAuth("",m.username,m.password,m.host)
	return nil
}

