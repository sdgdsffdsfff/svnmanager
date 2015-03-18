package task

import (
	"time"
	"king/helper"
	"king/bootstrap"
)


type Task struct {
	IsRunning bool
	Enable bool
	Watch func()
	Duration time.Duration
}

func (r *Task) Stop(){
	r.Enable = false
}

func (r *Task) Start() {
	r.Enable = true
}

func (r *Task) Quit(){
	r.IsRunning = false
}

var defaultDuration = time.Second * 5
var taskList = map[string]*Task{}
var taskLoopEnabled = false

func New(name string, callback func(*Task), duration ...time.Duration){
	d := defaultDuration

	if len(duration) > 0 {
		d = duration[0]
	}

	task := &Task{
		IsRunning: false,
		Enable: false,
		Duration: d,
	}

	task.Watch = func(){
		if task.IsRunning {
			return
		}
		task.IsRunning = true
		for {
			if task.IsRunning == false {
				break
			}
			if taskLoopEnabled && task.Enable {
				callback(task)
			}
			time.Sleep(task.Duration)
		}
	}

	taskList[name] = task
}

func Trigger(name string){

	StartTask()

	if method, found:= taskList[name]; found {
		method.Enable = true
		method.IsRunning = true
	}
}


func StartTask(){
	taskLoopEnabled = true
}
func StopTask(){
	taskLoopEnabled = false
}

func init(){
	bootstrap.Register(func(){
		StartTask()
		helper.AsyncMap(taskList, func(key, value interface{}) bool {
			value.(*Task).Watch()
			return false
		})
	})
}
