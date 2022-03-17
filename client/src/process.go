package ofisclient

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/process"
)

func ListProcess() ([]*process.Process, error) {
	return process.Processes()
}

func CheckProcess(name string) bool {
	list, err := ListProcess()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}

	for _, a := range list {
		neym, err := a.Name()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		fmt.Printf("process: %v\n", neym)
		if neym == name {
			return true
		}
	}

	return false
}
