package logi

import (
	"fmt"
	"os"
)

func Log(lines ...any) error {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		_, err := f.WriteString(fmt.Sprint(line) + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func Clear() error {
	return os.WriteFile("log.txt", []byte{}, 0644)
}
