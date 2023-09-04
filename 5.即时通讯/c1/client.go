package c1

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"sync"
	"talk/util"
	"time"
)

const (
	goCount      = 2
	sendChanSize = 5
)

type ClientInterface interface {
	ReadFromWS()
	Write2WS()
	Run()
}

func NewClientInterface(conn *websocket.Conn, userID int) ClientInterface {
	return NewClient(conn, userID)
}

func NewClient(con *websocket.Conn, userID int) *Client {
	return &Client{userID, con, make(chan any, sendChanSize), goCount, &sync.WaitGroup{}}
}

type Client struct {
	userID int
	*websocket.Conn
	sendChan  chan any
	goCount   int
	waitGroup *sync.WaitGroup
}

type MessageFromWS struct {
	FromUserID   int       `json:"fromUserID"`
	FromUserName string    `json:"fromUserName"`
	ToUserID     int       `json:"toUserID"`
	ToUserName   string    `json:"toUserName"`
	Payload      string    `json:"payload"`
	SendTime     time.Time `json:"sendTime"`
	IsBroadcast  bool      `json:"isBroadcast"`
	ToGroupID    int       `json:"toGroupID"`
	ToGroupName  string    `json:"toGroupName"`
}

func (c *Client) ReadFromWS() {
	defer func() {
		//撤销从全局map
		MessageChan <- MessageUnRegister{c}
		c.waitGroup.Done()
	}()
	//注册至全局map
	MessageChan <- MessageRegister{c}
	var message MessageFromWS
	var err error
	var b []byte
	for {
		_, b, err = c.Conn.ReadMessage()
		if err != nil {
			//如果客户端主动关闭,也会触发错误
			util.Logger.Errorf("ReadPump ReadMessage err:%v", err)
			return
		}
		err = json.Unmarshal(b, &message)
		if err != nil {
			util.Logger.Errorf("ReadPump Unmarshal err:%v", err)
			return
		}
		message.SendTime = time.Now()
		util.Logger.Debugf("message:%#v", message)
		if message.IsBroadcast {
			MessageChan <- MessageBroadcast{MessageBase{
				FromUserID:   message.FromUserID,
				FromUsername: message.FromUserName,
				Payload:      message.Payload,
				SendTime:     message.SendTime,
			}, message.ToGroupID, message.ToGroupName}
		} else {
			MessageChan <- Message{MessageBase{
				FromUserID:   message.FromUserID,
				FromUsername: message.FromUserName,
				Payload:      message.Payload,
				SendTime:     message.SendTime,
			}, message.ToUserID, message.ToUserName}
		}
	}
}

func (c *Client) Write2WS() {
	defer func() {
		c.waitGroup.Done()
	}()
	var message any
	var err error
	var b []byte
	var ok bool
	var w io.WriteCloser
	for {
		message, ok = <-c.sendChan
		if !ok {
			err = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			if err != nil {
				util.Logger.Errorf("Write2WS WriteMessage CloseMessage err:%v", err)
				return
			}
		}

		b, err = json.Marshal(message)
		if err != nil {
			util.Logger.Errorf("Write2WS Marshal err:%v", err)
		}
		w, err = c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			util.Logger.Errorf("Write2WS NextWriter err:%v", err)
			return
		}
		w.Write(b)
		n := len(c.sendChan)
		for i := 0; i < n; i++ {
			w.Write([]byte("\n"))
			b, err = json.Marshal(<-c.sendChan)
			w.Write(b)
		}

	}
}

func (c *Client) Run() {
	c.waitGroup.Add(goCount)
	go c.ReadFromWS()
	go c.Write2WS()
	c.waitGroup.Wait()
}
