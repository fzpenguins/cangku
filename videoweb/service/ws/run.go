package ws

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"videoweb/database/DB/dao"
	"videoweb/database/DB/model"
	"videoweb/database/cache"
	"videoweb/pkg/e"
	"videoweb/rabbitmq"
)

func MessageHandler() {
	for {
		select {
		case broadcast := <-Broadcast:
			//log.Println("建立新连接")

			msgDao := dao.GetMsgDao(context.Background())

			if broadcast.Type == e.SingleChat { //单聊
				msgJSON, _ := json.Marshal(broadcast)
				//存储到数据库

				if !Clients[ClientsSingleMap[broadcast.FromUid]] { //对方有无上线
					err := rabbitmq.PublishMsg(msgJSON, broadcast.Target)
					if err != nil {
						log.Println(err)
						return
					}
					log.Println("对方未上线")
					err = msgDao.StoreSingleChatMsg(broadcast.FromUid, broadcast.Target, broadcast.Content, false)
					if err != nil {
						log.Println(err)
						return
					}
					continue
				}
				select {
				case ClientsSingleMap[broadcast.FromUid].send <- msgJSON: //发送
					err := msgDao.StoreSingleChatMsg(broadcast.FromUid, broadcast.Target, broadcast.Content, true)
					if err != nil {
						log.Println(err)
						return
					}
				default:
					close(ClientsSingleMap[broadcast.FromUid].send)

					delete(ClientsSingleMap, broadcast.FromUid)
				}

			} else if broadcast.Type == e.GroupChat { //群聊

				var members = make([]string, e.MaxStore)
				members, err := cache.RedisClient.SMembers(context.Background(), broadcast.Target).Result()
				if err != nil {
					log.Println(err)
				}
				msgJSON, _ := json.Marshal(broadcast)
				onlineMember := ClientsMap[broadcast.Target] //这里的问题只存入一个，为什么会这样
				for _, member := range members {

					if member == broadcast.FromUid {
						continue
					}

					msg := &model.Message{
						CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
						DeletedAt: "",
						FromUid:   member,
						ToUid:     broadcast.Target,
						Type:      e.GroupChat,
						Content:   broadcast.Content,
					}
					//log.Println(onlineMember)
					if conn, ok := onlineMember[member]; ok {
						conn.send <- msgJSON
						msg.ReadTag = true
					} else {
						//log.Println(member)
						msg.ReadTag = false
					}

					if err := dao.Db.Where(&model.Message{}).Create(&msg).Error; err != nil {
						log.Println(err)
						continue
					}
				}
			}
		case conn := <-Unregister:
			log.Println("结束连接")
			ClientsLock.Lock()

			if conn.Type == e.SingleChat {
				if Clients[conn] {
					delete(ClientsSingleMap, conn.FromUid)
					delete(Clients, conn)
					close(conn.send)
				}
			} else if conn.Type == e.GroupChat {
				delete(Clients, conn)
				delete(ClientsMap[conn.Target], conn.FromUid)
				close(conn.send)
			}
			if len(ClientsMap[conn.Target]) == 0 {
				delete(ClientsMap, conn.Target)
			}

			ClientsLock.Unlock()
			log.Println("已完全结束连接")
		case conn := <-Register:
			log.Println("进行register相关操作")

			if conn.Type == e.GroupChat { //type 还没赋值，无法进入此分支
				log.Println(conn.Target)
				//log.Println(conn.FromUid)
				if _, ok := ClientsMap[conn.Target]; !ok {
					ClientsMap[conn.Target] = make(map[string]*Client) //不能反复make，有make的就不能再次make了
				}

				ClientsLock.Lock()
				Clients[conn] = true
				ClientsMap[conn.Target][conn.FromUid] = conn
				ClientsLock.Unlock()

				cache.RedisClient.SAdd(context.Background(), conn.Target, conn.FromUid)
			} else if conn.Type == e.SingleChat {

				ClientsLock.Lock()
				Clients[conn] = true
				ClientsSingleMap[conn.Target] = conn //目标映射Client
				ClientsLock.Unlock()

			}

		}
	}
}
