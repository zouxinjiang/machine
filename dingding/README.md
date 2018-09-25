# 钉钉机器人
## Variable
```go
//消息类型
const (
	TEXT MsgType = iota     //text类型
	LINK                    //link类型
	MARKDOWN                //markdown类型
	SINGLE_ACTIONCARD       //整体跳转ActionCard类型
	ISLAND_ACTIONCARD       //独立跳转ActionCard类型
	FEEDCARD                //FeedCard类型
)
```
## Types
```go
//机器人结构
type  Machine struct {
	prefixVar string   //变量前缀
	url string          //机器人地址
	method string       //提交数据的方法。默认POST
	msg string          //发送的消息
	msgtype machine.MsgType //发送的消息类型
	tmpl map[string]string  //消息模板
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
```

## API调用顺序
1. m := dingding.NewMachine(机器人地址)
2. 生成模板消息。m.GenerateTmpl(消息类型,map[string]interface{}{该消息类型所需要的数据})
3. 将模板注册，方便以后使用。 m.RegisteTmpl(模板名字,生成的模板消息)
4. 使用模板构建具体消息。m.BuildTxtMsg(之前注册的模板名, map[string]string{自己定义的模板需要的数据})
5. 发送消息。m.SendMsg()

[注意]模板默认变量前缀为“$”,如果要改变变量前缀，请在BuildTxtMsg()之前调用SetPrefixVar()
## 程序示例
```go
package main

import (
	"machine/dingding"
	"fmt"
)

func main() {
    //创建机器人
	m := dingding.NewMachine("https://oapi.dingtalk.com/xxx")
	//生成一个模板
	tmpl,err := m.GenerateTmpl(dingding.ISLAND_ACTIONCARD, map[string]interface{}{"TITLE":"t111","TEXT":"text","BTNS":[]dingding.Btn{dingding.Btn{"1","2"},dingding.Btn{"11","22"}}})
	if err != nil {
		panic(err)
	}
	//注册模板
	m.RegisteTmpl("textmsg",tmpl)
	//使用模板构建消息
	m.BuildTxtMsg("textmsg", map[string]string{"USER":"xxx"})
    //发送消息
	fmt.Println(m.SendMsg())
}
```
## 消息类型以及字段
参考官方文档：https://open-doc.dingtalk.com/microapp/serverapi2/qf2nxq

类型 | 字段 | 类型 | 说明
---|---|---|---
TEXT | content | string | 消息内容
LINK | title | string | 消息标题
-| text | string | 消息内容。如果太长只会部分展示
-| messageUrl | string | 点击消息跳转的URL
-| picUrl | string | 图片URL
MARKDOWN| title | string | 首屏会话透出的展示内容
-| text | string | markdown格式的消息
SINGLE_ACTIONCARD | title | string | 首屏会话透出的展示内容
-| text | string | markdown格式的消息
-| singleTitle | string | 单个按钮的方案。(设置此项和singleURL后btns无效)
-| singleURL | string | 点击singleTitle按钮触发的URL
ISLAND_ACTIONCARD | title | string | 首屏会话透出的展示内容
-| text | string | markdown格式的消息
-| btns | json数组字符串 | 按钮的信息：title-按钮方案，actionURL-点击按钮触发的URL
FEEDCARD | links | json数组字符串 | 条目
links | title | string | 单条信息文本
links | messageURL | string | 点击单条信息到跳转链接
links | picURL | string | 单条信息后面图片的URL


