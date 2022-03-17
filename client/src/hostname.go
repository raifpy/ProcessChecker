package ofisclient

import "github.com/shirou/gopsutil/v3/host"

func Hostname() (string, error) {

	inf, err := host.Info()
	if err != nil {
		return "", err
	}

	return inf.Hostname, nil

}
