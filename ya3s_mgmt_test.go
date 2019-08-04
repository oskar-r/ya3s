package ya3s

import (
	"log"
	"testing"
)

func TestListRegisteredTasks(t *testing.T) {
	Setup()
	AddTask(testFunc, "* * * *")
	AddTask(testFunc2, "05 02 * *")
	tests := []struct {
		name string
	}{
		{
			"TEST_1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ListRegisteredTasks()
		})
	}
	for {
	}
}

func testFunc() error {
	log.Printf("TestFunc")
	return nil
}

func testFunc2() error {
	log.Printf("TestFunc")
	return nil
}
