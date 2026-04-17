package helper

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteToFile(offset int64, bytes *[]byte, filePath string) {

	// 1. Extract the directory path from the full file path
	dir := filepath.Dir(filePath)

	// 2. Create the directory tree (0755 is standard permissions)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}

	defer f.Close()

	_, _ = f.WriteAt(*bytes, offset)
}
