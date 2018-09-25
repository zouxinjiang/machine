# 邮箱机器人
## Variable
```go
//消息类型
const (
	TEXT MsgType = iota     //text类型
)
```
## Types
```go
//机器人结构
type Machine struct {
	auth smtp.Auth      //认证对象
	msg string             //消息
	msgType machine.MsgType //消息类型
	username string         //发送者用户名
	password string         //密码
	host string             //邮箱验证地址/发送邮件地址【不加端口】
	prefixVar string        //模板变量前缀。默认$
	tmpl map[string]string  //模板
	port string             //端口。默认25
	to []string             //发送到
}
```

## API调用顺序
1. m := mail.NewMachine(机器人地址)
2. 登录。 m.Login(登录信息)
3. 生成模板消息。m.GenerateTmpl(消息类型,map[string]interface{}{该消息类型所需要的数据})
4. 将模板注册，方便以后使用。 m.RegisteTmpl(模板名字,生成的模板消息)
5. 使用模板构建具体消息。m.BuildTxtMsg(之前注册的模板名, map[string]string{自己定义的模板需要的数据})
6. 发送消息。m.SendMsg()

[注意]模板默认变量前缀为“$”,如果要改变变量前缀，请在BuildTxtMsg()之前调用SetPrefixVar()
## 程序示例
```go
package main

import (
	"machine/mail"
	"fmt"
)

func main() {
	m := mail.NewMachine("smtp.qq.com")
	m.Login("xxx@qq.com","xxx","smtp.qq.com")
	tpl,err := m.GenerateTmpl(mail.TEXT,map[string]interface{}{"subject":"主题","content":"模板内容$CC","to":[]string{"yyyy@qq.com","zzz@163.com"}})
	if err != nil {
		panic(err)
	}
	m.RegisteTmpl("textTmpl",tpl)
	m.BuildTxtMsg("textTmpl",map[string]string{"CC":"替换模板"})
	fmt.Println(m.SendMsg())

```
## 消息类型以及字段
类型 | 字段 | 类型 | 说明
---|---|---|---
TEXT | subject | string | 邮件主题
-| content | string | 邮件内容
-| to | []string | 邮件发送给谁



