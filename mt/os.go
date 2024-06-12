package mt

import (
	"os"
	"path/filepath"

	"github.com/shirou/gopsutil/v3/process"
)

func IsInRun(pros string) bool {
	psList, _ := process.Processes()
	for _, ps := range psList {
		name, _ := ps.Name()
		if pros == name {
			return true
		}
	}

	return false
}

var CurrentDir string

func GetCurrentDir() string {
	if CurrentDir != "" {
		return CurrentDir
	}

	ex, err := os.Executable()
	if err != nil {
		return ""
	}

	CurrentDir = filepath.Dir(ex)
	return CurrentDir
}

func GetJoinPathFromHere(paths ...string) string {
	currentdir := GetCurrentDir()
	newPaths := []string{currentdir}
	newPaths = append(newPaths, paths...)
	return filepath.Join(newPaths...)
}
