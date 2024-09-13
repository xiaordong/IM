package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GroupChatMessage 群聊消息结构
type GroupChatMessage struct {
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	GroupID   uint   `json:"group_id"` // 群组ID
}

// GroupChatClient 群聊客户端
type GroupChatClient struct {
	ID      string
	GroupID uint
	Socket  *websocket.Conn
	Send    chan []byte `json:"-"`
}

// GroupChatManager 群聊管理器
type GroupChatManager struct {
	Clients   map[uint][]*GroupChatClient // 以群组ID为键，存储该群组的所有客户端
	Broadcast chan *GroupChatMessage      // 广播消息通道
	Online    chan *GroupChatClient       // 客户端上线事件
	Offline   chan *GroupChatClient       // 客户端下线事件
}

var GroupManager = GroupChatManager{
	Clients:   make(map[uint][]*GroupChatClient),
	Broadcast: make(chan *GroupChatMessage),
	Online:    make(chan *GroupChatClient),
	Offline:   make(chan *GroupChatClient),
}

// Start 启动群聊服务
func (manager *GroupChatManager) Start() {
	for {
		log.Println("--- 监听群聊管道通信 ---")
		select {
		case client := <-manager.Online:
			log.Printf("客户端上线: %v\n", client.ID)
			manager.Clients[client.GroupID] = append(manager.Clients[client.GroupID], client)
			log.Printf("客户端 %s 加入群组: %d\n", client.ID, client.GroupID)
		case client := <-manager.Offline:
			log.Printf("客户端下线: %v\n", client.ID)
			manager.removeClient(client)
			log.Printf("客户端 %s 已从群组: %d 移除\n", client.ID, client.GroupID)
		case message := <-manager.Broadcast:
			BroadcastGroupChatMessages(message)
		}
	}
}

// removeClient 从管理器中移除客户端
func (manager *GroupChatManager) removeClient(client *GroupChatClient) {
	log.Printf("正在移除客户端: %s\n", client.ID)
	for i, c := range manager.Clients[client.GroupID] {
		if c.ID == client.ID {
			manager.Clients[client.GroupID] = append(manager.Clients[client.GroupID][:i], manager.Clients[client.GroupID][i+1:]...)
			log.Printf("客户端: %s 已被移除\n", client.ID)
			break
		}
	}
}

// BroadcastGroupChatMessages 广播消息给群组内的所有客户端
func BroadcastGroupChatMessages(message *GroupChatMessage) {
	clients, exists := GroupManager.Clients[message.GroupID]
	if !exists {
		log.Println("没有找到群组内的客户端")
		return
	}
	log.Printf("群组: %d 有 %d 个客户端在线\n", message.GroupID, len(clients))
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Println("消息序列化失败:", err)
		return
	}
	for _, client := range clients {
		if client.ID == message.Sender {
			continue
		}
		select {
		case client.Send <- msgBytes:
			log.Printf("消息发送给客户端: %s\n", client.ID)
		default:
			log.Printf("客户端: %s 消息发送通道已满，丢弃消息\n", client.ID)
		}
	}
}

// GroupWsHandler 处理WebSocket连接请求
func GroupWsHandler(c *gin.Context) {
	groupIDStr := c.Query("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的群组ID"})
		return
	}
	uid := c.Query("uid")

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法升级到WebSocket"})
		return
	}

	client := &GroupChatClient{
		ID:      uid,
		GroupID: uint(groupID),
		Socket:  conn,
		Send:    make(chan []byte, 100), // 设置缓冲大小
	}

	log.Printf("客户端: %s 连接成功，加入群组: %d\n", uid, groupID)
	GroupManager.Online <- client
	go client.Read()
	go client.Write()
}

// StartGroupChatService 启动群聊服务
func StartGroupChatService() {
	log.Println("群聊服务启动...")
	go GroupManager.Start()
}

// Read 从WebSocket读取消息
func (c *GroupChatClient) Read() {
	defer func() {
		GroupManager.Offline <- c
		_ = c.Socket.Close()
		log.Printf("客户端: %s 已下线\n", c.ID)
	}()

	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("客户端: %s 正常关闭连接\n", c.ID)
			} else {
				log.Printf("客户端: %s 读取消息失败，错误: %v\n", c.ID, err)
			}
			break
		}
		var receivedMsg struct {
			Sender  string `json:"sender"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal(msg, &receivedMsg); err != nil {
			log.Printf("客户端: %s 反序列化消息失败: %v\n", c.ID, err)
			continue
		}
		groupMessage := &GroupChatMessage{
			Sender:    receivedMsg.Sender,
			Content:   receivedMsg.Content,
			Timestamp: time.Now().Unix(), // 生成时间戳
			GroupID:   c.GroupID,
		}

		// 广播消息
		GroupManager.Broadcast <- groupMessage
	}
}

// Write 向WebSocket写入消息
func (c *GroupChatClient) Write() {
	defer func() {
		_ = c.Socket.Close()
		log.Printf("客户端: %s 的写入协程结束\n", c.ID)
	}()

	ticker := time.NewTicker(time.Second * 10) // 定期发送心跳
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				log.Printf("客户端: %s 的消息发送通道已关闭\n", c.ID)
				return
			}
			log.Printf("客户端: %s 正在发送消息: %s\n", c.ID, message)
			if err := c.Socket.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("客户端: %s 发送消息失败: %v\n", c.ID, err)
				return
			}
		case <-ticker.C:
			// 发送心跳或其他保活消息
			if err := c.Socket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("客户端: %s 发送心跳失败: %v\n", c.ID, err)
				return
			}
			log.Printf("客户端: %s 发送心跳成功\n", c.ID)
		}
	}
}
