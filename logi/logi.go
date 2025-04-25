package logi

import (
	"fmt"
	"os"
)

func Log(line any) error {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprint(line) + "\n")
	return err
}

func Clear() error {
	return os.WriteFile("log.txt", []byte{}, 0644)
}
