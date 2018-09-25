package dingding

import (
	"net/http"
	"strings"
	"machine"
	"encoding/json"
)

//实现了Machine下的MsgType接口
type MsgType int

const (
	TEXT MsgType = iota
	LINK
	MARKDOWN
	SINGLE_ACTIONCARD
	ISLAND_ACTIONCARD
	FEEDCARD
)

func (mt MsgType)ToString() string{
	return mt.String()
}

func (mt MsgType)String() string{
	switch mt {
	case TEXT:
		return "TEXT"
	case LINK:
		return "LINK"
	case MARKDOWN:
		return "MARKDOWN"
	case ISLAND_ACTIONCARD:
		return "ISLAND_ACTIONCARD"
	case SINGLE_ACTIONCARD:
		return "SINGLE_ACTIONCARD"
	case FEEDCARD:
		return "FEEDCARD"
	default:
		return "unkown msg type"
	}
}

//预设模板
type tmpl string
const (
	textTmpl tmpl = `
{
    "msgtype": "text", 
    "text": {
        "content": "${CONTENT}$"
    }
}
`
	linkTmpl tmpl = `{
    "msgtype": "link", 
    "link": {
        "text": "${TEXT}$", 
        "title": "${TITLE}$", 
        "picUrl": "${PICURL}$", 
        "messageUrl": "${MESSAGEURL}$"
    }
}`
	markdownTmpl tmpl = `
{
     "msgtype": "markdown",
     "markdown": {
         "title":"${TITLE}$",
         "text": "${TEXT}$"
     }
 }`
	singleActionCardTmpl tmpl =`
{
	"msgtype": "actionCard",
    "actionCard": {
        "title": "${TITLE}$", 
        "text": "${TEXT}$", 
        "hideAvatar": "0", 
        "btnOrientation": "0", 
        "singleTitle" :"${SINGLETITLE}$",
        "singleURL" : "${SINGLEURL}$"
    }
}`
	islandActionCardTmpl tmpl = `
{
	"msgtype": "actionCard",
    "actionCard": {
        "title":"${TITLE}$", 
        "text": "${TEXT}$", 
        "hideAvatar": "0", 
        "btnOrientation": "0", 
        "btns": ${BTNS}$
    }
}`
	feedCardTmpl tmpl = `
{
 	"msgtype": "feedCard",
    "feedCard": {
        "links": ${LINKS}$
    }	
}`
)


type  Machine struct {
	prefixVar string
	url string
	method string
	msg string
	msgtype machine.MsgType
	tmpl map[string]string
}
//feedCard类型使用的link
type Link struct{
	Title string `json:"title"`
	MessageUrl string `json:"messageURL"`
	PicUrl string `json:"picURL"`
}
//islandActionCard类型使用的btn
type Btn struct {
	Title string `json:"title"`
	ActionUrl string `json:"actionURL"`
}


func NewMachine(url string) *Machine{
	return &Machine{
		url:url,
		method:"post",
		tmpl:make(map[string]string),
		msgtype:TEXT,
		prefixVar:"$",
	}
}
func (m Machine)Login(data map[string]string) error{
	return nil
}
func (m *Machine)SetPrefixVar(prefix string){
	m.prefixVar = prefix
}

func (m Machine) SendMsg() (int, error) {
	l := len(m.msg)
	_,err := http.Post(m.url,"application/json",strings.NewReader(m.msg))
	if err != nil {
		return -1,machine.NewSendErr("发送数据失败，请检查网络")
	}
	return l,err
}

func (m *Machine) RegisteTmpl(name,tmpl string) error {
	m.tmpl[name] = tmpl
	return nil
}

//指定模板生成消息
func (m *Machine) BuildTxtMsg(tmplName string, data map[string]string) error {
	str,ok := m.tmpl[tmplName]
	if !ok {
		return machine.NewTmplErr("模板"+tmplName+"不存在")
	}
	for k,v := range data {
		str = strings.Replace(str,m.prefixVar+k,v,-1)
	}
	m.msg = str
	return nil
}

func (m Machine) BuildStreamMsg(tmplName string, data interface{}) error {
	panic("implement me")
}


//生成模板
func (m Machine) GenerateTmpl(msgType machine.MsgType,data map[string]interface{}) (tpl string,curerr error){
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

	m.msgtype = msgType

	switch msgType {
	case TEXT:
		tpl = strings.Replace(string(textTmpl),"${CONTENT}$",data["CONTENT"].(string),-1)
	case LINK:
		tpl = strings.Replace(string(linkTmpl),"${TEXT}$",data["TEXT"].(string),-1)
		tpl = strings.Replace(tpl,"${TITLE}$",data["TITLE"].(string),-1)
		tpl = strings.Replace(tpl,"${PICURL}$",data["PICURL"].(string),-1)
		tpl = strings.Replace(tpl,"${MESSAGEURL}$",data["MESSAGEURL"].(string),-1)
	case MARKDOWN:
		tpl = strings.Replace(string(markdownTmpl),"${TEXT}$",data["TEXT"].(string),-1)
		tpl = strings.Replace(tpl,"${TITLE}$",data["TITLE"].(string),-1)
	case SINGLE_ACTIONCARD:
		tpl = strings.Replace(string(singleActionCardTmpl),"${TEXT}$",data["TEXT"].(string),-1)
		tpl = strings.Replace(tpl,"${TITLE}$",data["TITLE"].(string),-1)
		tpl = strings.Replace(tpl,"${SINGLETITLE}$",data["SINGLETITLE"].(string),-1)
		tpl = strings.Replace(tpl,"${SINGLEURL}$",data["SINGLEURL"].(string),-1)
	case ISLAND_ACTIONCARD:
		tpl = strings.Replace(string(islandActionCardTmpl),"${TEXT}$",data["TEXT"].(string),-1)
		tpl = strings.Replace(tpl,"${TITLE}$",data["TITLE"].(string),-1)
		btns,ok := data["BTNS"]
		if !ok {
			curerr = machine.NewTmplErr("生成模板出错。类型"+msgType.ToString()+",缺少BTNS参数")
		}
		var xx interface{}
		res,ok := btns.([]map[string]string)
		res2,ok2 := btns.([]Btn)

		if !ok && !ok2 {
			curerr = machine.NewTmplErr("生成模板出错。类型"+msgType.ToString()+",BTNS参数类型错误。BTNS类型应该为[]map[string]string或[]Btn")
			break
		}else if ok {
			xx = res
		}else {
			xx = res2
		}

		btnsStr,err := json.Marshal(xx)
		curerr = err
		tpl = strings.Replace(tpl,"${BTNS}$",string(btnsStr),-1)
	case FEEDCARD:
		lnks,ok := data["LINKS"]
		if !ok {
			curerr = machine.NewTmplErr("生成模板出错。类型"+msgType.ToString()+",缺少LINKS参数")
		}
		var xx interface{}
		res,ok := lnks.([]map[string]string)
		res2,ok2 := lnks.([]Link)
		if !ok && !ok2 {
			curerr = machine.NewTmplErr("生成模板出错。类型"+msgType.ToString()+",LINKS参数类型错误。LINKS类型应该为[]map[string]string或[]Link")
			break
		}else if ok {
			xx = res
		} else {
			xx = res2
		}
		links,err := json.Marshal(xx)
		curerr = err
		tpl = strings.Replace(string(feedCardTmpl),"${LINKS}$",string(links),-1)
	}
	return tpl,curerr
}



