package pkg
import (
	"fmt"
	"os"
	"path/filepath"
)

func CheckFilePath(path string) (string,error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return "",err
	}
	return absPath,nil
}
func OpenFile(path string) (*os.File, error) {
	absPath, err := CheckFilePath(path)
	if err!=nil {
		return nil,err
	}
	file, err := os.Open(absPath)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Incorrect file path or file does not exist")
		return nil, err
	}
	return file, nil
}
