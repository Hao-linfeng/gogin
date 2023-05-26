package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

// 消息
type Message struct {
	gorm.Model
	FormId   int64  //发送者
	TargetId int64  //接收者
	Type     int    //发送类型
	Media    int    //消息类型
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var ClientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// 发送者ID 接收者ID 消息类型 发送内容 发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	//校验
	// token := query.Get("token")
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	// msgType := query.Get("type")
	// targetId := query.Get("targetId")
	// context := query.Get("context")
	isvalida := true //checkToken
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//获取token
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//用户关系

	// userid 跟 node绑定 并加锁
	rwLocker.Lock()
	ClientMap[userId] = node
	rwLocker.Unlock()

	//发送逻辑
	go sendprocess(node)

	//接收逻辑
	go recvProc(node)

	sendMsg(userId, []byte("欢迎进入聊天室"))

}

func sendprocess(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws] <<<<<<<", data)
	}
}

var upsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	upsendChan <- data
}

func init() {
	go udpSendProc()
	go udpPrevProc()
}

// 完成udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 201),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {

		select {
		case data := <-upsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据接收协程
func udpPrevProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println("接收error>>>>>>>>>", err)
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		sendMsg(msg.TargetId, data)
		// case 2://群发
		// 	sendGroupMsg()
		// case 3://广播
		// 	sendAllMsg()
		// case 4:

	}

}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := ClientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg

	}

}
