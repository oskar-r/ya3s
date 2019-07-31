/*
yasss - Yet Another Super Simple Scheduler is a task schedule tkat executed tasks on assigned interval based on a cron syntax

At the begning of each minute an assessmen if any task is up for execution, if so the task is executed and reports back execution statis to the task que. Tasks need to contain all logic needed to execute them and should only return an error
*/
package ya3s

import (
	"testing"
	"time"
)

func Test_timeToExecute(t *testing.T) {
	type args struct {
		interval string
		t        time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"TEST_1_EXEC",
			args{
				interval: "* * * *",
				t:        time.Date(2019, 06, 16, 13, 57, 04, 0, time.UTC),
			},
			true,
		},
		{
			"TEST_2_EXEC",
			args{
				interval: "57 * * *",
				t:        time.Date(2019, 06, 16, 13, 57, 04, 0, time.UTC),
			},
			true,
		},
		{
			"TEST_3_EXEC",
			args{
				interval: "57 13 15 *",
				t:        time.Date(2019, 06, 15, 13, 57, 04, 0, time.UTC),
			},
			true,
		},
		{
			"TEST_4_NO_EXEC",
			args{
				interval: "* * 15 *",
				t:        time.Date(2019, 06, 16, 13, 57, 04, 0, time.UTC),
			},
			false,
		},
		{
			"TEST_4_EXEC_MULTI",
			args{
				interval: "5,10,15,20 * 16 *",
				t:        time.Date(2019, 06, 16, 13, 10, 04, 0, time.UTC),
			},
			true,
		},
		{
			"TEST_5_EXEC_MULTI",
			args{
				interval: "0,10,20 * * *",
				t:        time.Date(2019, 06, 16, 13, 00, 04, 0, time.UTC),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeToExecute(tt.args.interval, tt.args.t); got != tt.want {
				t.Errorf("timeToExecute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateSchedule(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"OK_SCHEMA",
			args{
				s: "* * * *",
			},
			false,
		},
		{
			"OK_SCHEMA_2",
			args{
				s: "5,10,15 * * *",
			},
			false,
		},
		{
			"OK_SCHEMA_3",
			args{
				s: "5,10,15 10 * *",
			},
			false,
		},
		{
			"BAD_SCHEMA_1",
			args{
				s: "5,10,15 A * *",
			},
			true,
		},
		{
			"BAD_SCHEMA_2",
			args{
				s: "5,10,15 *  * *",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateSchedule(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("validateSchedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
