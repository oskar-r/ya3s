# Ya3s - Yet another super simplistic scheduler
This is a work in progress to provide a simple way of scheduling stand alone functions to run at a specific point in time.

## Example
```go
import (
   "fmt"
   "github.com/oskarr/ya3s"
)

func testFunc() error {
    //execute a task
    fmt.Printf("Task executed")
    return nil
}
func testFunc2() error {
    //execute a task
    fmt.Printf("Task executed")
    return nil
}

func main() {
    //Set up the task scheduler
    ya3s.Setup()
    //Add task to run every minute
    ya3s.AddTask(testFunc, "* * * *") 
    //Add task to run every tenth minute
    ya3s.AddTask(testFunc2, "0,10,20,30,40,50 * * *")
}
```
At this point in time you must your self secure that the task can exxecute within the timeframe for the task

## Todo
* Dependencies between tasks scheduled to run
* REST interface to query the task scheduler for status and running tasks