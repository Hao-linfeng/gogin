package models

import (
	"fmt"
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
	FormId   string //发送者
	TargetId string //接收者
	Type     string //发送类型
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
	msgType := query.Get("type")
	targetId := query.Get("targetId")
	context := query.Get("context")
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

}
