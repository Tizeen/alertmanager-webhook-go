# Alertmanager Webhook

接收Alertmanager发送的请求，通过[Pushbear](http://pushbear.ftqq.com/admin/#/)将告警消息推送到微信，并记录日志。

感谢[Easy](https://github.com/easychen)提供这么好的服务。

##使用

1. 注册Pushbear，得到通道的sendKey；
1. 填入`sendKey`变量中；
1. 编译：`go build -o alertwebhook main.go`；
1. 赋予执行权限，直接运行；