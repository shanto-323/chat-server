package ui

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	APP_NAME string = "Chat App"
)

type App struct {
	app         fyne.App
	userListBox *fyne.Container
	socket      *Socket
}

func NewApp() *App {
	a := app.New()
	socket := NewSocket()
	return &App{
		app:         a,
		userListBox: container.NewVBox(),
		socket:      socket,
	}
}

func (u *App) Run() {
	w := u.app.NewWindow(APP_NAME)
	w.SetContent(u.Ui())
	w.ShowAndRun()
}

func (u *App) Ui() fyne.CanvasObject {
	topBar := u.topBar()
	listBar := container.NewVScroll(u.userListBox)
	chatBar := container.NewVBox()

	r1 := canvas.NewRectangle(color.White)
	r2 := canvas.NewRectangle(color.White)

	objects := []fyne.CanvasObject{topBar, listBar, chatBar, r1, r2}
	ui := container.New(NewUiLayout(topBar, listBar, chatBar, r1, r2), objects...)
	return ui
}

func (u *App) topBar() fyne.CanvasObject {
	button := widget.NewButton("Connect", nil)

	button.OnTapped = func() {
		if u.socket != nil && u.socket.Connected {
			if err := u.socket.Disconnect(); err != nil {
				return
			}
			u.userListBox.Objects = nil
			u.userListBox.Refresh()
			button.SetText("Connect")
		} else {
			err := u.socket.Connect()
			if err != nil {
				button.SetText("Connect")
			} else {
				go func() {
					ticker := time.NewTicker(2 * time.Second)
					defer ticker.Stop()
					for range ticker.C {
						if u.socket.Connected == false {
							u.userListBox.Objects = nil
							u.userListBox.Refresh()
							button.SetText("Connect")
							break
						}
					}
				}()
				go u.writeList()
				button.SetText("Disconnect")
			}
		}
	}
	return button
}

func (u *App) writeList() {
	for {
		fmt.Println("running")
		select {
		case <-u.socket.Ctx.Done():
			return
		case msg := <-u.socket.NewList:
			if msg {
				u.userListBox.Objects = nil
				for _, id := range u.socket.Clients.IdList {
					u.userListBox.Add(widget.NewButton(id, func() {
						log.Println(id)
					}))
				}
				u.userListBox.Refresh()
			}
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
				SenderId:   u.socket.Client.Id,
				ReceiverId: idEntity.Text,
				Payload:    chatEntity.Text,
			}

			if err := u.socket.conn.WriteJSON(message); err != nil {
				log.Println(err)
				return
			}
		},
	)

	messageSection := container.NewVBox(chatEntity, idEntity, button)
	return messageSection
}
