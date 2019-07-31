package ya3s

import (
	"log"
	"testing"
)

func TestListRegisteredTasks(t *testing.T) {
	Setup()
	AddTask(testFunc, "* * * *")

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
}

func testFunc() error {
	log.Printf("TestFunc")
	return nil
}
