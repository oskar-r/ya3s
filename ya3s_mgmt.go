package ya3s

import (
	"log"
	"reflect"
	"regexp"
	"runtime"
)

//ListRegisteredTasks returns task and status for all currently registered tasks
func ListRegisteredTasks() {
	for k, v := range tm.tasks {
		var re = regexp.MustCompile(`(?m).*\.(.*)$`)
		tn := runtime.FuncForPC(reflect.ValueOf(v.task).Pointer()).Name()
		log.Printf("%+v", tn)
		t := re.FindSubmatch([]byte(tn))

		if len(t) > 1 && t[1] != nil {
			tn = string(t[1])
		}
		log.Printf("TaskID:%s\tname:%s\t", k, tn)
	}
}
