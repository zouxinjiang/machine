package machine

/*
调用顺序:
创建一个机器人 NewMachine()
生成模板 GenerateTmpl()
注册模板 RegisteTmpl()
[设置模板变量前缀 SetPrefixVar()]
构建消息 BuildStreamMsg()
发送消息 SendMsg()
 */

//机器人接口
type Machine interface {
	//发送消息
	SendMsg() (int,error)
	//注册消息模板
	RegisteTmpl(name,tmpl string) error
	//构建文本消息
	BuildTxtMsg(tmplName string,data map[string]string) error
	//构建二进制消息
	BuildStreamMsg(tmplName string,data interface{}) error
	//设置变量前缀
	SetPrefixVar(prefix string)
	//生成模板信息
	GenerateTmpl(msgType MsgType,data map[string]interface{}) (tpl string,curerr error)
	//登录机器人
	Login(data map[string]string) error

}


