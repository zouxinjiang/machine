# 机器人发送消息
## 目前支持机器人
机器人平台 | 是否支持
---|---
钉钉自定义机器人 | 是

## 使用方法
1. m := NewMachine(机器人地址)
2. 登录机器人。 m.Login(登录信息)
3. 生成模板消息。m.GenerateTmpl(消息类型,map[string]interface{}{该消息类型所需要的数据})
4. 将模板注册，方便以后使用。 m.RegisteTmpl(模板名字,生成的模板消息)
5. 使用模板构建具体消息。m.BuildTxtMsg(之前注册的模板名, map[string]string{自己定义的模板需要的数据})
6. 发送消息。m.SendMsg()




