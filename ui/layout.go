package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

var WINDOW_SIZE fyne.Size = fyne.NewSize(500, 250)

type UiLayout struct {
	topBar  fyne.CanvasObject
	listBar fyne.CanvasObject
	chatBar fyne.CanvasObject
	r1      *canvas.Rectangle
	r2      *canvas.Rectangle
}

func NewUiLayout(
	topBar, listBar, chatBar fyne.CanvasObject,
	r1, r2 *canvas.Rectangle,
) fyne.Layout {
	return &UiLayout{
		topBar:  topBar,
		listBar: listBar,
		chatBar: chatBar,
		r1:      r1,
		r2:      r2,
	}
}

func (u *UiLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	totalHeight := size.Height
	totalWidth := size.Width

	u.topBar.Resize(fyne.NewSize(totalWidth, totalHeight*0.10))
	u.topBar.Move(fyne.NewPos(0, 0))

	u.listBar.Resize(fyne.NewSize(totalWidth*0.30, totalHeight*0.90))
	u.listBar.Move(fyne.NewPos(0, totalHeight*0.10))

	u.chatBar.Resize(fyne.NewSize(totalWidth*0.70, totalHeight*0.90))
	u.chatBar.Move(fyne.NewPos(totalWidth*0.30, totalHeight*0.10))

	u.r1.Resize(fyne.NewSize(totalWidth, 1))
	u.r1.Move(fyne.NewPos(0, totalHeight*0.10))
	u.r2.Resize(fyne.NewSize(1, totalHeight*0.90))
	u.r2.Move(fyne.NewPos(totalWidth*0.30, totalHeight*0.10))
}

func (u *UiLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return WINDOW_SIZE
}
