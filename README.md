## 简介
go-wol是一个使用Go编写的支持 WOL（Wake-on-LAN）功能的程序，支持命令行接口（CLI）和 Web API 两种方式。可以通过命令行或 HTTP 请求触发唤醒目标设备。

## WOL 工作原理
WOL（Wake-on-LAN）是通过发送一个特定的“魔术包”（Magic Packet）来唤醒网络中处于睡眠或关机状态的计算机。

## 运行

### API 模式 

不加任何参数，直接运行程序，默认监听8080端口。

```bash
go run main.go
```
然后，发送 POST 请求到 `http://localhost:8080/wakeonlan`，请求体为：

```json
{
  "mac_address": "00-14-22-01-23-45"
}
```
`MAC`地址支持`:`和`-`分隔符格式，可以使用 curl 命令行工具来测试：

```bash
curl -X POST -H "Content-Type: application/json" -d '{"mac_address": "00:14:22:01:23:45"}' http://localhost:8080/wakeonlan
```

### 命令行模式

指定 `-mac` 参数将启动 CLI 模式： `-mac` 参数后输入MAC地址来向改地址发送魔术包。

```bash

go run main.go -mac "00:11:22:33:44:55"
```
## 编译

运行项目根目录下的 `build.sh`，可以编译出适用于不同架构、操作系统的可执行文件，脱离Go环境运行。
```bash
bash build.sh
```