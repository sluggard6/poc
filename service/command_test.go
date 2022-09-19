package service

import (
	"fmt"
	"testing"
)

func TestServiceImpl(t *testing.T) {
	cs := GetCommandService()
	stdout, stderr, err := cs.Run("ls")
	if err == nil {
		fmt.Println("------------------")
		fmt.Print(stdout)
		fmt.Print(stderr)
	} else {
		fmt.Println(err.Error())
	}
}
