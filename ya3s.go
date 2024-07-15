/*
Package ya3s - Yet Another Super Simple Scheduler is a task schedule that executed tasks on assigned interval based on a cron syntax

At the beginning of each minute assess if any task is up for execution, if so the task is executed and reports back execution statistics to the task que. Tasks need to contain all logic needed to execute them and should only return an error

Copyright (C) 2019 by Oskar Roman <roman.oskar@gmail.com>
Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:
//
The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.
//
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package ya3s

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

const (
	errWrongNrOfItems   = "schedule must have exactly 4 items separated by blank-space or tab found %d"
	errScheduleInterval = "execution interval %s is not * or between (%d-%d)"
	errScheduleFmt      = "format of schedule %s not correct"
)

// Task is a function that executes when scheduled
type Task func() error

type taskItem struct {
	task                    Task
	lastExecution           time.Time
	lastExecutionSuccessful bool
	lastError               error
	numberOfExecutions      int64
	schedule                string
}

type taskMap struct {
	tasks map[string]*taskItem
}

var tm taskMap

// Setup a new execution que and start up the execution clock. Should be called before AddTask
func Setup() {
	tm.tasks = make(map[string]*taskItem)
	ticker := time.NewTicker(time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				assessTaskMap() //assessQue looks in to execution que to see if any task is up for execution and if so executes this task.
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// CleanUp is a function that sets the task queue to nil
func CleanUp() {
	fmt.Printf("Cleaning up befor quiting...\n")
	tm.tasks = nil
}

// AddTask add a new task to the que for execution once in the defined schedule (cron style * * * *)
func AddTask(task Task, schedule string) (string, error) {
	schedule = strings.Replace(schedule, "\t", " ", -1) //Normalize to spaces
	if err := validateSchedule(schedule); err != nil {
		return "", err
	}
	myUuid := uuid.New()
	tm.tasks[myUuid.String()] = &taskItem{
		task:     task,
		schedule: schedule,
	}
	return myUuid.String(), nil
}

func assessTaskMap() {
	//*,*,*,*
	t := time.Now()
	for k, v := range tm.tasks {
		if timeToExecute(v.schedule, t) {
			go func(i *taskItem, k string, t time.Time) {
				err := i.task()
				if err != nil {
					tm.tasks[k].lastError = err
					tm.tasks[k].lastExecutionSuccessful = false
				} else {
					tm.tasks[k].lastError = nil
					tm.tasks[k].lastExecutionSuccessful = true
				}
				tm.tasks[k].numberOfExecutions++
				tm.tasks[k].lastExecution = t

			}(v, k, t)
		}
	}
}

// Validation of provided schedule
func validateSchedule(s string) error {
	si := strings.Split(s, " ")
	if len(si) != 4 {
		return fmt.Errorf(errWrongNrOfItems, len(si))
	}
	//Validate minute
	if err := validateScheduleItem(si[0], 0, 59); err != nil {
		return fmt.Errorf(err.Error(), "minute")
	}
	if err := validateScheduleItem(si[1], 0, 23); err != nil {
		return fmt.Errorf(err.Error(), "hour")
	}
	if err := validateScheduleItem(si[2], 1, 31); err != nil {
		return fmt.Errorf(err.Error(), "day")
	}
	if err := validateScheduleItem(si[3], 1, 7); err != nil {
		return fmt.Errorf(err.Error(), "weekday")
	}

	return nil
}

func validateScheduleItem(s string, low, high int) error {
	if s == "*" {
		return nil
	}
	si := strings.Split(s, ",")
	for _, v := range si {
		err := validPointInTime(v, low, high)
		if err != nil {
			return err
		}
	}
	return nil
}

func validPointInTime(s string, low, high int) error {
	i, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf(errScheduleFmt, s)
	}
	if i < low || i > high {
		return fmt.Errorf(errScheduleInterval, "%s", low, high)
	}
	return nil
}

// Assess if task should be executed based on the tasks schedule
func timeToExecute(schedule string, t time.Time) bool {
	iv := strings.Split(schedule, " ")

	//Weekday (replace with 0 if Sunday is scheduled as 7)
	if !execute(strings.Replace(iv[3], "7", "0", 1), int(t.Weekday())) {
		return false
	}
	//Day
	if !execute(iv[2], t.Day()) {
		return false
	}
	//Hour
	if !execute(iv[1], t.Hour()) {
		return false
	}
	//Minute
	if !execute(iv[0], t.Minute()) {
		return false
	}
	return true
}

func execute(i string, t int) bool {
	if i == "*" {
		return true
	}
	si := strings.Split(i, ",")
	for _, v := range si {
		et, err := strconv.Atoi(v)
		if err == nil && et == t {
			return true
		}
	}
	return false
}
