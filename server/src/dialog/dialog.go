package dialog

import (
	"github.com/gen2brain/dlgs"
)

type BuildinInfo struct {
	AppName string
}

var dialoginfo = BuildinInfo{
	AppName: "Server-Client based process name based confirmation application",
}

func SeedDialogInfo(b BuildinInfo) {
	dialoginfo = b
}

func FatalError(message string) {

	/*if usegui {
		dialog.Message("%s\nÖlümcül hata: %s\nÇıkış yapılacak.", dialoginfo.AppName, message).Error()
	} else {
		log.Printf("%s\n%s", dialoginfo.AppName, message)
	}*/

	dlgs.Error("Ölümcül Hata - "+dialoginfo.AppName, message)
}
