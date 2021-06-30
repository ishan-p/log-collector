package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ishan-p/log-collector/internal/schema"
)

type FsStorageConfig struct {
	BaseDir string `json:"base_dir"`
}

func getSubDir() string {
	now := time.Now()
	subDir := filepath.Join(fmt.Sprint(now.Year()), fmt.Sprint(now.Month()), fmt.Sprint(now.Day()), fmt.Sprint(now.Hour()))
	return subDir
}

func createDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0740)
		if err != nil {
			log.Println("Could not create directory: ", err)
			return err
		}
	}
	return nil
}

func writeFs(logEvent schema.CollectCmdPayload, baseDir string) {
	subDir := getSubDir()
	dir := filepath.Join(baseDir, subDir)
	createDirIfNotExists(dir)

	jsonLogEvent, err := json.Marshal(logEvent)
	if err != nil {
		log.Println("Unable to encode log as json")
	}

	fileName := "collector.log"
	appendLog(filepath.Join(dir, fileName), string(jsonLogEvent))
}

func appendLog(filePath string, data string) {
	// TODO: Add file lock
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		fmt.Println(err)
	}
	n, err := io.WriteString(f, data+"\n")
	if err != nil {
		fmt.Println(n, err)
	}
}
