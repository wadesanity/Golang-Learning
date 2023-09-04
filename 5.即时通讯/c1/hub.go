package c1

import (
	"talk/util"
	"time"
)

var (
	MessageChan        = make(chan any)
	registeredMap      = make(map[int]*Client)
	registeredGroupMap = make(map[int][]int)
)

type Message struct {
	MessageBase
	ToUserID   int
	ToUserName string
}

type MessageBroadcast struct {
	MessageBase
	ToGroupID   int
	ToGroupName string
}

type MessageBase struct {
	FromUserID   int
	FromUsername string
	Payload      string
	SendTime     time.Time
}

type MessageRegister struct {
	ClientAddr *Client
}

type MessageUnRegister struct {
	ClientAddr *Client
}

func InitHub() {
	go func() {
		var originMessage any
		var toClient *Client
		var ok bool
		var groupIDList = make([]int, 0)
		groupIDList = append(groupIDList, []int{1, 2}...)
		registeredGroupMap[1] = groupIDList
		for {
			originMessage = <-MessageChan
			util.Logger.Debugf("originMessage:%#v", originMessage)
			util.Logger.Debugf("map:%#v", registeredMap)
			switch message := originMessage.(type) {
			case MessageBroadcast:
				groupIDList, ok = registeredGroupMap[message.ToGroupID]
				if ok {
					for _, clientID := range groupIDList {
						if toClient, ok = registeredMap[clientID]; ok {
							select {
							case toClient.sendChan <- message:
							default: //存在阻塞,认为客户端挂了
								delete(registeredMap, clientID)
								close(toClient.sendChan)
							}
						}

					}
				}

			case Message:
				toClient, ok = registeredMap[message.ToUserID]
				if ok {
					select {
					case toClient.sendChan <- message:
					default: //存在阻塞,认为客户端挂了
						delete(registeredMap, toClient.userID)
						close(toClient.sendChan)
					}
				}
				//todo 当peer不在线
			case MessageRegister:
				registeredMap[message.ClientAddr.userID] = message.ClientAddr
			case MessageUnRegister:
				delete(registeredMap, message.ClientAddr.userID)
				close(message.ClientAddr.sendChan)
			default:
				util.Logger.Errorf("message not invalid:%#v", message)
			}
		}
	}()
}
