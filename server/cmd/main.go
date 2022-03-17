package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	ofis "server/src"
	"server/src/dialog"
)

const VERSION = 0.1
const SYSTEM = "SERVER"

func main() {
	f, err := os.Open("serverconfig.json")
	if err != nil {
		dialog.FatalError(err.Error())
		os.Exit(1)
	}
	defer f.Close()

	var cfg ofis.Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		dialog.FatalError(err.Error())
		os.Exit(1)
	}

	fmt.Printf("VERSION: %v\n", VERSION)
	fmt.Printf("SYSTEM: %v\n", SYSTEM)
	build, _ := debug.ReadBuildInfo()
	fmt.Printf("BUILD: %s\n", build.Path)
	fmt.Printf("SERVE ADDR %s\n", cfg.ListenAddr)

	dialog.SeedDialogInfo(dialog.BuildinInfo{
		AppName: "ProcessChecker Server",
	})

	software, err := ofis.NewOfis(cfg)

	if err != nil {
		dialog.FatalError(err.Error())
		os.Exit(1)
	}

	if err := software.Listen(); err != nil {
		dialog.FatalError(err.Error())
	}

	log.Println("PROGRAM FINISHED")
}
