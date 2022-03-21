package main

import (
	ofisclient "client/src"
	"client/src/dialog"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/adrg/xdg"
)

const APPNAME = "ProcessChecker Client"

var runapp = flag.String("runapp", "", "run application with controll")

func init() {
	flag.Parse()
}

func main() {

	var cfg ofisclient.Config

	path, err := xdg.ConfigFile("ofisclient.json")
	if err != nil {
		dialog.FatalError(err.Error())
		os.Exit(1)
	}
	fmt.Printf("path: %v\n", path)

	f, err := os.Open(path)
	if err != nil {
		dialog.FatalError(err.Error())
		os.Exit(1)
	}

	defer f.Close()

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		dialog.FatalError(err.Error())
		os.Exit(1)
	}

	dialog.SeedDialogInfo(dialog.BuildinInfo{
		AppName: APPNAME,
	})

	if *runapp != "" {

		var process string
		split := strings.Split(*runapp, string(os.PathSeparator))
		if len(split) > 1 {
			process = split[len(split)-1]
		} else {
			process = *runapp
		}
		log.Println("Kontrollü program çalıştırılacak!")
		dialog.PopUp(fmt.Sprintf("%s açılmadan önce ağda kontrol edilecek. Lütfen bekleyin", process))
		up, err := url.Parse(cfg.DialAddr)
		if err != nil {
			dialog.FatalError(err.Error())
			os.Exit(1)
		}
		hostname, err := ofisclient.Hostname()
		if err != nil {
			dialog.FatalError(err.Error())
			os.Exit(1)
		}
		response, err := http.Get("http://" + up.Host + "/sorgu?process=" + url.QueryEscape(process) + "&hostname=" + url.QueryEscape(hostname))
		if err != nil {
			dialog.FatalError(err.Error())
			os.Exit(1)
		}
		defer response.Body.Close()

		var list = []ofisclient.Sorgu{}
		if err := json.NewDecoder(response.Body).Decode(&list); err != nil {
			dialog.FatalError(err.Error())
			os.Exit(1)
		}
		var isopened bool
		var openedhost string
		for _, abc := range list {
			if abc.Process {
				isopened = true
				openedhost = abc.Hostname
				break
			}
		}

		if isopened {
			if !dialog.Soru(fmt.Sprintf("Program %s tarafından kullanılıyor!\n\nYine de çalıştırmak ister misiniz?", openedhost)) {
				os.Exit(0)
			}
		}
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmdlist := []string{"/C", "/Q", "start", *runapp}
			cmdlist = append(cmdlist, flag.Args()...)
			cmd = exec.Command("cmd.exe", cmdlist...)
		} else {
			cmd = exec.Command(*runapp, flag.Args()...)
		}
		if err := cmd.Start(); err != nil {
			dialog.FatalError(err.Error())
		}

		os.Exit(0)

	}

	ofis, err := ofisclient.NewOfisClient(cfg)
	if err != nil {
		dialog.FatalError(err.Error())
		os.Exit(1)
	}

	dialog.PopUp(fmt.Sprintf("İyi günler %s. %s aktif.", ofis.Hostname, APPNAME))

	if err := ofis.Run(); err != nil {
		dialog.FatalError(err.Error())
		os.Exit(1)
	}
}
