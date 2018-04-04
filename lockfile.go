package lockfile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/process"
)

// lock will look for the .lock file named after the main program and check
// among the running PIDs whether it's still running.
// It returns the new locking file's name and any errors.
func Lock() (string, error) {
	fileBaseName := strings.TrimRight(filepath.Base(os.Args[0]),
		filepath.Ext(os.Args[0]))
	lockFile := fmt.Sprintf("./%s.lock", fileBaseName)
	currentPid := []byte(strconv.Itoa(os.Getpid()))

	if _, err := os.Stat(lockFile); err == nil {
		pid, err := ioutil.ReadFile(lockFile)
		if err != nil {
			return "", fmt.Errorf("can't read the lockfile for the PID\n%s\n",
				err.Error())
		}

		pids, err := process.Pids()
		if err != nil {
			return "", fmt.Errorf("can't get list of PIDs, %s\n", err.Error())
		}

		pidInt64, err := strconv.ParseInt(string(pid), 10, 32)
		if err != nil {
			return "", fmt.Errorf("%s\n", err.Error())
		}
		pidInt32 := int32(pidInt64)

		pidExists := false
		for _, v := range pids {
			if v == pidInt32 {
				pidExists = true
				break
			}
		}

		if pidExists {
			return "", errors.New("another instance is locking this run")
		}
	}

	ioutil.WriteFile(lockFile, currentPid, 0600)

	return lockFile, nil
}
