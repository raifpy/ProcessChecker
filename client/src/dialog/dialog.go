package dialog

import (
	"fmt"

	"github.com/gen2brain/beeep"
	"github.com/gen2brain/dlgs"
)

type BuildinInfo struct {
	AppName string
}

var dialoginfo = BuildinInfo{
	AppName: "OFIS CLIENT",
}

func SeedDialogInfo(b BuildinInfo) {
	dialoginfo = b
}

func FatalError(message string) {
	fmt.Printf("FATAL: %v\n", message)
	dlgs.Error("FATAL "+dialoginfo.AppName, message)
}
func Error(message string) {
	dlgs.Error(dialoginfo.AppName, message)

}

func Info(message string) {
	dlgs.Info(dialoginfo.AppName, message)

}

func Soru(message string) bool {
	b, _ := dlgs.Question(dialoginfo.AppName, message, false)
	return b

}

func PopUp(message string) {
	beeep.Notify(dialoginfo.AppName, message, "")
}
