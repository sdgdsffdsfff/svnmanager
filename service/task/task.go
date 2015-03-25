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
	Timer *time.Timer
	OnComplete func(*Task)
	Result interface{}
	lock chan int
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

func New(name string, callback func(*Task) interface{}, duration ...time.Duration) *Task{
	d := defaultDuration

	if len(duration) > 0 {
		d = duration[0]
	}

	task := &Task{
		IsRunning: false,
		Enable: false,
		OnComplete: func(*Task){},
		lock: make(chan int),
	}

	task.Watch = func(){
		if task.IsRunning {
			return
		}
		task.IsRunning = true
		for {
			if taskLoopEnabled && task.Enable {
				task.Result = nil
				task.Result = callback(task)
				task.OnComplete(task)
				go func(){
					time.Sleep( d )
					task.lock <- 1
				}()
			}
			<- task.lock
		}
	}

	taskList[name] = task

	return task
}

func Trigger(name string, fn ...func(*Task)){

	StartTask()

	if method, found:= taskList[name]; found {
		method.Enable = true
		method.IsRunning = true
		method.lock <- 1
		if len(fn) > 0 {
			method.OnComplete = fn[0]
		}
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
