package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"github.com/gorilla/websocket"
)

type Socket struct {
	conn      *websocket.Conn
	mu        sync.RWMutex
	Ctx       context.Context
	Cancel    context.CancelFunc
	Client    *UserModel
	Clients   UserList
	Connected bool
}

func NewSocket() *Socket {
	return &Socket{
		Clients: UserList{},
	}
}

func (s *Socket) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Println(err)
		return err
	}
	s.Ctx, s.Cancel = context.WithCancel(context.Background())
	s.conn = conn
	s.Connected = true

	go s.ReadMessage()
	go s.UpdateList()

	return nil
}

func (s *Socket) Disconnect() error {
	var err error
	if s.Cancel == nil || s.conn == nil {
		err = fmt.Errorf("")
		return err
	}

	s.Cancel()
	s.conn.Close()
	s.Connected = false
	return nil
}

func (s *Socket) ReadMessage() {
	for {
		select {
		case <-s.Ctx.Done():
			log.Println("stoping read-message")
			return
		default:
			{
				_, msg, err := s.conn.ReadMessage()
				if err != nil {
					log.Println("WebSocket read error:", err)
					s.Disconnect()
					return
				}

				var message Message
				if err := json.Unmarshal(msg, &message); err != nil {
					log.Println(err)
					continue
				}

				switch message.MsgType {
				case TYPE_INFO:
					var client UserModel
					if err := json.Unmarshal(message.Payload, &client); err != nil {
						log.Println(err)
						continue
					}
					s.Client = &client
				case TYPE_CHAT:
					log.Println(string(message.Payload))
				case TYPE_LIST:
					var list UserList
					if err := json.Unmarshal(message.Payload, &list); err != nil {
						log.Println(err)
						continue
					}
					if len(s.Clients.IdList) == len(list.IdList) {
						continue
					}
					s.mu.Lock()
					s.Clients = list
					s.mu.Unlock()

					log.Println(list)
					fyne.Do(
						func() {
							// u.WriteList()
						},
					)
				}
			}
		}
	}
}

func (s *Socket) UpdateList() {
	for {
		for {
			select {
			case <-s.Ctx.Done():
				log.Println("stoping update-list")
				return
			default:
				{
					if s.Client == nil {
						continue
					}
					s.mu.Lock()
					message := Message{
						MsgType:  TYPE_LIST,
						SenderId: s.Client.Id,
					}
					if err := s.conn.WriteJSON(&message); err != nil {
						log.Println(err)
						return
					}
					s.mu.Unlock()
					time.Sleep(5 * time.Second)
				}
			}
		}
	}
}
