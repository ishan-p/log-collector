package client

import (
	"os"

	"github.com/hpcloud/tail"
)

func getOffset(fileName string) tail.SeekInfo {
	fileInfo, _ := os.Stat(fileName)
	return tail.SeekInfo{
		Offset: fileInfo.Size(),
	}
	// TODO: Add state management logic to maintain last read state
}
