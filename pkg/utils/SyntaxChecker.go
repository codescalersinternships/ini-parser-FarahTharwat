package utils
import (
	"fmt"
	"os"
	"path/filepath"
	//"strings"
)

func OpenFile(filename string) (*os.File, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return nil, err
	}
	file, err := os.Open(absPath)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Incorrect file path or file does not exist")
		return nil, err
	}
	return file, nil
}