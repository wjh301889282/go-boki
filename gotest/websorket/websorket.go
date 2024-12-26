package websorket

import (
	"encoding/json"
	"exchangeapp/global"
	"exchangeapp/models/user"
	"exchangeapp/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

// PendingInvitation 存储待处理邀请信息
type PendingInvitation struct {
	InviterID int    `json:"inviter_id"`
	InviteeID int    `json:"invitee_id"`
	GroupID   int    `json:"group_id"`
	Message   string `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[int]*websocket.Conn) // 存储用户ID与WebSocket连接
var lastHeartbeat = make(map[int]time.Time) // 存储每个用户的最后心跳时间
var clientsMutex sync.Mutex

var heartbeatTimeout = 5 * time.Second // 如果超时超过60秒认为用户掉线

// WebSocket连接处理
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP请求到WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection:", err)
		}
	}(conn)

	// 获取用户 token（假设第一个消息是身份验证消息）
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}

	userID := 0
	// 解析消息
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err == nil {
		if msg["type"] == "authenticate" {
			token, ok := msg["token"].(string)
			if ok {
				// 验证 token
				userID, err = getUserID(token)
				if err != nil {
					log.Println("Token validation failed:", err)
					return
				}
				// 将用户连接加入连接池
				clientsMutex.Lock()
				clients[userID] = conn
				lastHeartbeat[userID] = time.Now() // 记录首次心跳时间
				clientsMutex.Unlock()
				log.Printf("用户 %d 认证成功", userID)
			} else {
				log.Println("No token in authentication message")
			}
		}
	}

	// 监听消息（包括心跳包）
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			// 锁定并删除该用户的连接
			clientsMutex.Lock()
			log.Printf("User %d is offline", userID)
			delete(clients, userID)
			clientsMutex.Unlock() // 解锁

			break
		}

		// 处理心跳包消息
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			if msg["type"] == "heartbeat" {
				// 更新该用户的心跳时间
				clientsMutex.Lock()
				lastHeartbeat[userID] = time.Now()
				clientsMutex.Unlock()
				//log.Printf("Received heartbeat from user %d", userID)
				continue // 如果是心跳包，直接跳过后面的处理
			}
		}

		// 处理接收到的消息，假设是邀请消息
		var invitation PendingInvitation
		err = json.Unmarshal(message, &invitation)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// 处理邀请
		handleInvitation(invitation)
	}
}

// 处理邀请逻辑
func handleInvitation(invitation PendingInvitation) {
	// 检查被邀请者是否在线
	clientsMutex.Lock()
	inviteeConn, online := clients[invitation.InviteeID]
	clientsMutex.Unlock()

	if online {
		// 被邀请者在线，发送邀请消息
		err := inviteeConn.WriteJSON(invitation)
		if err != nil {
			log.Println("Error sending invitation:", err)
		} else {
			fmt.Println("Invitation sent to the invitee")
		}
	} else {
		// 被邀请者不在线，存储邀请信息到数据库
		storePendingInvitation(invitation)
	}
}

// 存储待处理的邀请到数据库
func storePendingInvitation(invitation PendingInvitation) {
	// 假设我们已经建立了数据库连接 db
}

// 获取用户ID
func getUserID(token string) (int, error) {

	// 解析 Token 获取用户名
	username, err := utils.ParseJWT(token)
	if err != nil {
		log.Printf("Failed to parse token: %v", err)
		return 0, fmt.Errorf("failed to parse token: %v", err)
	}

	// 根据 username 查询用户
	var user user.User
	if err := global.Db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("User not found: %v", err)
		return 0, fmt.Errorf("user not found")
	}

	// 返回用户 ID
	return int(user.ID), nil
}

// 每隔一定时间检查所有用户的在线状态
func checkHeartbeat() {
	ticker := time.NewTicker(10 * time.Second) // 每10秒检查一次
	for {
		select {
		case <-ticker.C:
			clientsMutex.Lock()
			for userID, conn := range clients {
				// 检查每个用户的最后心跳时间
				if time.Since(lastHeartbeat[userID]) > heartbeatTimeout {
					// 如果超过超时，认为该用户掉线
					log.Printf("User %d is offline", userID)
					// 可以选择从连接池中移除，或者将其状态更新为离线
					delete(clients, userID)
					conn.Close()
				}
			}
			clientsMutex.Unlock()
		}
	}
}

// Initwebsorket 启动 WebSocket 服务器
func Initwebsorket() error {
	http.HandleFunc("/echo", handleWebSocket) // 设置处理 WebSocket 的路由
	fmt.Println("sorket 配置成功")

	//监听心跳包简易版
	//go checkHeartbeat()

	// 启动 WebSocket 服务并监听端口
	err := http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		return fmt.Errorf("WebSocket 服务启动失败: %v", err)
	}
	return nil
}
