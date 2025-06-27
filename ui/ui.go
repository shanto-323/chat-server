package ui

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/websocket"
	"github.com/tinrab/retry"
)

const (
	APP_NAME string = "Chat App"
)

var WINDOW_SIZE fyne.Size = fyne.NewSize(500, 250)

type App struct {
	app         fyne.App
	conn        *websocket.Conn
	client      *UserModel
	clients     UserList
	userListBox *fyne.Container
	mu          sync.RWMutex
}

func NewApp() *App {
	a := app.New()
	return &App{
		app:         a,
		clients:     UserList{},
		userListBox: container.NewVBox(),
	}
}

func (u *App) Run() {
	w := u.app.NewWindow(APP_NAME)
	w.Resize(WINDOW_SIZE)
	w.SetContent(u.CurrentScreen())
	w.ShowAndRun()
}

func (u *App) CurrentScreen() fyne.CanvasObject {
	appContainer := container.NewBorder(
		container.NewHBox(
			widget.NewButton("Connect", func() {
				go func() {
					retry.ForeverSleep(
						2*time.Second,
						func(_ int) error {
							conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
							if err != nil {
								log.Println(err)
								return err
							}
							u.conn = conn
							return nil
						},
					)
					go u.ReadMessage()
					go u.UpdateList()
				}()
			}),
		),
		nil,
		u.userListBox,
		nil,
		u.WriteMessage(),
	)
	return appContainer
}

func (u *App) ReadMessage() {
	for {
		_, msg, err := u.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
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
			u.client = &client
		case TYPE_CHAT:
			log.Println(string(message.Payload))
		case TYPE_LIST:
			var list UserList
			if err := json.Unmarshal(message.Payload, &list); err != nil {
				log.Println(err)
				continue
			}
			if len(u.clients.IdList) == len(list.IdList) {
				continue
			}
			u.mu.Lock()
			u.clients = list
			u.mu.Unlock()

			log.Println(list)
			fyne.Do(
				func() {
					u.WriteList()
				},
			)
		}
	}
}

func (u *App) WriteMessage() fyne.CanvasObject {
	chatEntity := widget.NewEntry()
	chatEntity.PlaceHolder = "Message"
	idEntity := widget.NewEntry()
	idEntity.PlaceHolder = "id"
	button := widget.NewButton(
		"send", func() {
			message := OutgoingMessage{
				MsgType:    TYPE_CHAT,
				SenderId:   u.client.Id,
				ReceiverId: idEntity.Text,
				Payload:    chatEntity.Text,
			}

			if err := u.conn.WriteJSON(message); err != nil {
				log.Println(err)
				return
			}
		},
	)

	messageSection := container.NewVBox(chatEntity, idEntity, button)
	return messageSection
}

func (u *App) WriteList() {
	u.userListBox.Objects = nil
	if len(u.clients.IdList) == 0 {
		return
	}

	for _, id := range u.clients.IdList {
		u.userListBox.Add(widget.NewButton(id, func() {
			log.Println(id)
		}))
	}
	u.userListBox.Refresh()
}

func (u *App) UpdateList() {
	for {
		if u.client == nil {
			continue
		}
		u.mu.Lock()
		message := Message{
			MsgType:  TYPE_LIST,
			SenderId: u.client.Id,
		}
		if err := u.conn.WriteJSON(&message); err != nil {
			log.Println(err)
			return
		}
		u.mu.Unlock()
		time.Sleep(5 * time.Second)
	}
}
