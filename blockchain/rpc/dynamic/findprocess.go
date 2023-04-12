package dynamic

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/duality-solutions/web-bridge/internal/util"
	"github.com/shirou/gopsutil/process"
)

func findProcess() (*process.Process, error) {
	processList, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("process.Processes() Failed")
	}
	for _, process := range processList {
		name, _ := process.Name()
		if strings.HasPrefix(name, DefaultName) {
			cmd, _ := process.Cmdline()
			// TODO check for same datadir as well
			if strings.Index(cmd, "-port="+strconv.Itoa(int(DefaultPort))) > 0 && strings.Index(cmd, "-rpcport="+strconv.Itoa(int(DefaultRPCPort))) > 0 {
				return process, nil
			}
		}
	}
	return nil, fmt.Errorf("findProcess Failed")
}

// FindDynamicdProcess returns the dynamicd processes running locally
func FindDynamicdProcess(start bool, timeout time.Duration) (*Dynamicd, *process.Process, error) {
	var dynamicd *Dynamicd
	var err error
	if start {
		dynamicd, err = LoadRPCDynamicd()
	}
	process, err := findProcess()
	if err == nil {
		return dynamicd, process, nil
	}
	for {
		select {
		case <-time.After(time.Second * 2):
			process, err := findProcess()
			if err == nil {
				return dynamicd, process, nil
			}
			if start {
				dynamicd, err = LoadRPCDynamicd()
				if err != nil {
					util.Error.Println("FindDynamicdProcess error restarting dynamicd process", err)
					continue
				}
			}
		case <-time.After(timeout):
			return dynamicd, nil, fmt.Errorf("Running dynamicd process not found")
		}
	}
}
