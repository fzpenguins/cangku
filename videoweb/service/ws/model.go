package ws

import (
	"encoding/json"
	"sync"
	"videoweb/database/DB/model"

	"github.com/hertz-contrib/websocket"
)

var ClientsLock sync.RWMutex

type Message struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
	PageNum int    `json:"page_num"`
}

type MsgFromBroadcast struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
	Target  string `json:"target"`
	FromUid string `json:"from_uid"`

	Flag bool `json:"flag"` //未读消息
}

//已读未读通过是否回复消息获得

// 定义客户端结构
type Client struct {
	conn    *websocket.Conn
	FromUid string
	send    chan []byte
	Type    int    //类型
	Target  string //目标id
	//groupId string
	//toUid string
}

type ClientManager struct {
	Clients    interface{}
	BroadCast  chan *MsgFromBroadcast
	Register   chan *Client
	Unregister chan *Client
}

// 创建新的客户端
func newClient(conn *websocket.Conn, uid string) *Client {
	return &Client{
		conn:    conn,
		send:    make(chan []byte, 1024*4),
		FromUid: uid,
		//Type: t,
		//Target: target,
	}
}

func (c *Client) SendHistory(msgs []*model.Message) {
	var err error
	for _, item := range msgs {
		buf := make([]byte, 0)

		if buf, err = json.Marshal(item); err == nil {

			c.conn.WriteMessage(websocket.TextMessage, buf)
		}

	}

}
