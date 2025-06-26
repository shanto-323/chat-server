package ui

import (
	"encoding/json"
	"log"
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

var WINDOW_SIZE fyne.Size = fyne.NewSize(900, 700)

type App struct {
	app     fyne.App
	conn    *websocket.Conn
	clients UserList
	table   *widget.Table
}

func NewApp() *App {
	a := app.New()
	return &App{
		app:     a,
		clients: UserList{},
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
		u.WriteList(),
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
		case TYPE_CHAT:
			log.Println(string(message.Payload))
		case TYPE_LIST:
			var list UserList
			if err := json.Unmarshal(message.Payload, &list); err != nil {
				log.Println(err)
				continue
			}
			u.clients = list
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
				MsgType:    TYPE_LIST,
				SenderId:   idEntity.Text,
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

func (u *App) WriteList() fyne.CanvasObject {
	log.Println("working")
	u.table = widget.NewTable(
		func() (int, int) {
			return len(u.clients.IdList), 1 // rows, columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell Template")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row < len(u.clients.IdList) {
				label.SetText(u.clients.IdList[id.Row])
			}
		},
	)

	return u.table
}

func (u *App) UpdateList() {
	for {
		if len(u.clients.IdList) != 0 {
			continue
		}

		time.Sleep(2 * time.Second)
		if u.table != nil {
			u.table.Refresh()
		}
	}
}
