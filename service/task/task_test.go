package task

import (
	"fmt"
	"os"
	"testing"
)

func Test_Task(t *testing.T) {
	th := NewTaskHandler()
	taskFunc := func() {
		fmt.Println("task test....")
	}
	if err := th.AddByFunc("1", "*/1 * * * * ?", taskFunc); err != nil {
		fmt.Printf("error to add TaskHandler task: %s", err)
		os.Exit(-1)
	}
	th.Start()
	select {}
}
