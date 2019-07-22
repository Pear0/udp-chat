package main

import (
	"fmt"
	"github.com/Pear0/udp-chat/ptypes"
	"github.com/jroimartin/gocui"
	"log"
	"strings"
	"time"
)

func guiMain(a *App) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layoutFactory(a))

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}


	go func() {
		for msg := range a.recvMessages {
			msg2 := msg
			g.Update(func(g *gocui.Gui) error {

				history, err := g.View("chat history")
				if err != nil {
					return err
				}

				_, _ = fmt.Fprintf(history, "%s: %s\n", msg2.SenderName, msg2.Message)

				return nil
			})
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func handleCommand(g *gocui.Gui, a *App, bufferText string) error {
	history, err := g.View("chat history")
	if err != nil {
		return err
	}

	switch {
	case strings.HasPrefix(bufferText, "/name"):
		a.Name = strings.TrimSpace(strings.TrimPrefix(bufferText, "/name"))
		_, _ = fmt.Fprintf(history, "[internal]: name is now: %s\n", a.Name)

	default:
		_, _ = fmt.Fprintf(history, "[internal]: unknown command: %s\n", bufferText)
		return nil
	}

	return nil
}

func handleMessage(g *gocui.Gui, a *App, bufferText string) error {
	if strings.HasPrefix(bufferText, "/") {
		return handleCommand(g, a, bufferText)
	}

	return a.Send(&ptypes.BasicMessage{
		SenderName: a.Name,
		Message: bufferText,
		Timestamp: uint32(time.Now().Unix()),
	})
}

func layoutFactory(a *App) func(*gocui.Gui) error {
	return func(g *gocui.Gui) error {
		maxX, maxY := g.Size()

		if v, err := g.SetView("chat history", 0, 0, maxX-1, maxY-3); err != nil {
			if err != gocui.ErrUnknownView || v == nil {
				return err
			}
			_, _ = fmt.Fprintln(v, "Hello world!")

			v.Autoscroll = true
		}

		if v, err := g.SetView("chat input", 0, maxY-3, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView || v == nil {
				return err
			}

			_, err := g.SetCurrentView("chat input")
			if err != nil {
				return err
			}

			err = g.SetKeybinding("chat input", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {

				buf := strings.TrimSpace(v.Buffer())
				if len(buf) == 0 {
					return nil
				}
				v.Clear()
				_ = v.SetCursor(0, 0)

				return handleMessage(g, a, buf)
			})
			if err != nil {
				return err
			}

			// v.Editor = Editor

			v.Highlight = true

			// v.FgColor = gocui.Attribute(15 + 1)
			// v.BgColor = gocui.Attribute(0)
			// v.BgColor = gocui.ColorDefault

			v.Autoscroll = false
			v.Editable = true
			v.Wrap = false
			v.Frame = false
		}

		return nil
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
