package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// 1. 定义命令行参数
	chain := flag.String("chain", "BSC", "要初始化的区块链 (例如: BSC, ETH)")
	name := flag.String("name", "My-CLI-Wallet", "为钱包自定义的名称")
	phone := flag.String("phone", "18888888888", "用户的手机号 (可选)")
	email := flag.String("email", "test@example.com", "用户的邮箱地址 (可选)")
	flag.Parse()

	// 2. 定义目标 API 地址
	url := "http://localhost:8888/api/wallet_init"

	// 3. 准备请求数据
	requestData := map[string]interface{}{
		"chain":        *chain,
		"name":         *name,
		"phone_number": *phone,
		"email":        *email,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Fatalf("错误: 无法打包 JSON 数据: %v", err)
	}

	// 4. 创建并发送 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("错误: 无法创建请求: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	fmt.Printf("正向 %s 发送请求...\n", url)
	fmt.Printf("请求体: %s\n", string(jsonData))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("错误: 发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 5. 读取并打印响应结果
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("错误: 读取响应体失败: %v", err)
	}

	fmt.Println("\n--- 响应结果 ---")
	fmt.Printf("HTTP 状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应体: %s\n", string(body))
}
