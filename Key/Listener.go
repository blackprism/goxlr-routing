package Key

import (
	"fmt"
	"github.com/blackprism/goxlr-routing/Config"
	"github.com/blackprism/goxlr-routing/GOXLR"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
)

type Listener struct {
	user32      *syscall.DLL
	keybindings []Config.Keybinding
	payload     GOXLR.Payload
}

func NewListener(actions []Config.Keybinding, payload GOXLR.Payload) Listener {
	return Listener{
		user32:      syscall.MustLoadDLL("user32"),
		keybindings: actions,
		payload:     payload,
	}
}

func (l *Listener) Register() {
	l.user32 = syscall.MustLoadDLL("user32")

	reghotkey := l.user32.MustFindProc("RegisterHotKey")
	for _, keybinding := range l.keybindings {
		keys := strings.Split(keybinding.Name, "+")
		mods := keys[:len(keys)-1]
		key := keys[len(keys)-1]

		modSum := 0
		for _, mod := range mods {
			switch mod {
			case "Ctrl":
				modSum += ModCtrl
			case "Alt":
				modSum += ModAlt
			case "Shift":
				modSum += ModShift
			}
		}

		// Register hotkeys:
		r1, _, err := reghotkey.Call(0, uintptr(keybinding.Id), uintptr(modSum), uintptr(key[0]))
		if r1 != 1 {
			fmt.Println("Failed to register", keybinding.Name, ", error:", err)
		}
	}
}

func (l Listener) Listen(c *websocket.Conn) {
	peekmsg := l.user32.MustFindProc("PeekMessageW")

	for {
		var msg = &MSG{}
		peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		// Registered id is in the WPARAM field:
		if id := msg.WPARAM; id != 0 {
			for _, keybinding := range l.keybindings {
				if keybinding.Id != id {
					continue
				}

				log.Println("Hotkey pressed:", keybinding.Name)

				for _, trigger := range keybinding.Triggers {
					if trigger.Type == "routing" {
						for _, action := range trigger.Actions {
							w, _ := c.NextWriter(websocket.TextMessage)
							w.Write(l.payload.Routing(action.Name, action.Input, action.Output))
							w.Close()
						}
					}
				}
			}
		}

		time.Sleep(time.Millisecond * 5)
	}
}
