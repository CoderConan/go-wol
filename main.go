package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// 发送魔术包的函数
func wakeOnLan(macAddress string) error {
	// 验证 MAC 地址格式
	if len(macAddress) != 17 {
		return fmt.Errorf("invalid MAC address format")
	}

	// 创建魔术包：前6个字节为 0xFF，后面重复6次目标 MAC 地址
	magicPacket := make([]byte, 102)

	// 填充魔术包的前 6 个字节为 0xFF
	for i := 0; i < 6; i++ {
		magicPacket[i] = 0xFF
	}

	// 提取 MAC 地址并重复 16 次
	macBytes := parseMacAddress(macAddress)
	for i := 6; i < len(magicPacket); i += len(macBytes) {
		copy(magicPacket[i:], macBytes)
	}

	// 设置目标广播地址和端口
	addr := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255), // 广播地址
		Port: 9,                            // WOL 默认使用端口 9
	}

	// 创建 UDP 连接
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 发送魔术包
	_, err = conn.Write(magicPacket)
	if err != nil {
		return err
	}

	// 给目标设备一点时间来接收包并唤醒
	time.Sleep(time.Second)

	return nil
}

// 解析 MAC 地址为字节数组，支持 `:` 和 `-` 分隔符
func parseMacAddress(macAddress string) []byte {
	// 去掉 MAC 地址中的分隔符（支持 `:` 和 `-`）
	macAddress = strings.ReplaceAll(macAddress, ":", "")
	macAddress = strings.ReplaceAll(macAddress, "-", "")

	// 如果 MAC 地址长度不等于 12 字符（6 字节），则无效
	if len(macAddress) != 12 {
		return nil
	}

	// 转换为字节数组
	macBytes := make([]byte, 6)
	for i := 0; i < 6; i++ {
		fmt.Sscanf(macAddress[i*2:i*2+2], "%x", &macBytes[i])
	}
	return macBytes
}

// CLI 入口
func runCLI(mac string) {
	if mac == "" {
		log.Println("MAC address is required")
		return
	}

	// 调用 WOL 函数
	log.Printf("sending magic packet to MAC address: %s\n", mac)
	if err := wakeOnLan(mac); err != nil {
		log.Printf("failed to send magic packet: %v\n", err)
	} else {
		log.Println("magic packet sent successfully!")
	}
}

// API 入口
func wakeOnLanHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MacAddress string `json:"mac_address"`
	}

	// 解析 JSON 请求体
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.MacAddress == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// 记录请求的 MAC 地址和请求者 IP 地址
	clientIP := r.RemoteAddr
	log.Printf("received request from IP: %s to wake up MAC address: %s\n", clientIP, req.MacAddress)

	// 调用 WOL 函数
	err = wakeOnLan(req.MacAddress)
	if err != nil {
		log.Printf("failed to send magic packet to %s: %v\n", req.MacAddress, err)
		http.Error(w, fmt.Sprintf("failed to send magic packet: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	log.Printf("magic packet successfully sent to %s\n", req.MacAddress)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("magic packet sent successfully!"))
}

func main() {
	// 定义命令行参数
	mac := flag.String("mac", "", "The MAC address to wake up (format: XX:XX:XX:XX:XX:XX or XX-XX-XX-XX-XX-XX)")
	flag.Parse()

	// 如果指定了 -mac 参数，执行 CLI 模式
	if *mac != "" {
		runCLI(*mac)
		return
	}

	// 启动 Web 服务器 API 模式
	http.HandleFunc("/wakeonlan", wakeOnLanHandler)

	// 启动 Web 服务器
	log.Println("starting WOL API server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start server: %v\n", err)
	}
}