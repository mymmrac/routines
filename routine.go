package routines

import (
	"fmt"
	"runtime"
	"time"
)

type Routine struct {
	started           bool
	completed         bool
	executionSequence []uintptr
	executed          map[uintptr]struct{}
	timers            map[uintptr]<-chan time.Time
	pc                []uintptr
}

func NewRoutine() *Routine {
	return &Routine{
		executionSequence: make([]uintptr, 0),
		executed:          make(map[uintptr]struct{}),
		timers:            make(map[uintptr]<-chan time.Time),
		pc:                make([]uintptr, 1),
	}
}

func (r *Routine) Started() bool {
	return r.started
}

func (r *Routine) Completed() bool {
	return r.completed
}

func (r *Routine) ExecutionSequenceRaw() []uintptr {
	executions := make([]uintptr, len(r.executionSequence))
	copy(executions, r.executionSequence)
	return executions
}

func (r *Routine) ExecutionSequence() []string {
	if len(r.executionSequence) == 0 {
		return nil
	}

	executions := make([]string, 0, len(r.executionSequence))
	frames := runtime.CallersFrames(r.executionSequence)
	for {
		frame, more := frames.Next()
		if frame.PC == 0 {
			continue
		}
		executions = append(executions, fmt.Sprintf("%s:%d", frame.File, frame.Line))

		if !more {
			break
		}
	}

	return executions
}

func (r *Routine) caller() uintptr {
	if runtime.Callers(3, r.pc) != 1 {
		panic(fmt.Errorf("failed to get caller"))
	}

	return r.pc[0]
}

func (r *Routine) isExecuted(caller uintptr) (executed bool) {
	_, executed = r.executed[caller]
	return executed
}

func (r *Routine) isPrevExecuted(myCaller uintptr) bool {
	for i, caller := range r.executionSequence {
		if i == len(r.executionSequence)-1 && caller == myCaller {
			return true
		}

		if !r.isExecuted(caller) {
			return false
		}
	}
	return true
}

func (r *Routine) addExecution(caller uintptr) {
	if len(r.executionSequence) > 0 && r.executionSequence[len(r.executionSequence)-1] == caller {
		return
	}

	r.executionSequence = append(r.executionSequence, caller)
}

func (r *Routine) markAsExecuted(caller uintptr) {
	r.executed[caller] = struct{}{}
}

func (r *Routine) executionTimer(caller uintptr, duration time.Duration) <-chan time.Time {
	if timer, found := r.timers[caller]; found {
		return timer
	}

	timer := time.After(duration)
	r.timers[caller] = timer
	r.addExecution(caller)
	return timer
}

func (r *Routine) Start() {
	if r.started {
		return
	}

	caller := r.caller()
	if r.isExecuted(caller) {
		return
	}
	if !r.isPrevExecuted(caller) {
		return
	}
	r.addExecution(caller)
	r.markAsExecuted(caller)

	r.started = true
}

func (r *Routine) End() {
	if !r.started {
		return
	}

	caller := r.caller()
	if r.isExecuted(caller) {
		return
	}
	if !r.isPrevExecuted(caller) {
		return
	}
	r.addExecution(caller)
	r.markAsExecuted(caller)

	r.started = false
	r.completed = true
}

func (r *Routine) Reset() {
	r.started = false
	r.completed = false
	r.executionSequence = make([]uintptr, 0)
	r.executed = make(map[uintptr]struct{})
	r.timers = make(map[uintptr]<-chan time.Time)
}

func (r *Routine) Restart() {
	r.Reset()
	r.Start()
}

func (r *Routine) Do(action func()) {
	if !r.started {
		return
	}

	caller := r.caller()
	if r.isExecuted(caller) {
		return
	}
	if !r.isPrevExecuted(caller) {
		return
	}
	r.addExecution(caller)
	r.markAsExecuted(caller)

	action()
}

func (r *Routine) WaitFor(duration time.Duration) {
	if !r.started {
		return
	}

	caller := r.caller()
	if r.isExecuted(caller) {
		return
	}
	if !r.isPrevExecuted(caller) {
		return
	}

	timer := r.executionTimer(caller, duration)
	select {
	case <-timer:
		r.markAsExecuted(caller)
	default:
		return
	}
}

func (r *Routine) WaitUntil(condition func() bool) {
	if !r.started {
		return
	}

	caller := r.caller()
	if r.isExecuted(caller) {
		return
	}
	if !r.isPrevExecuted(caller) {
		return
	}

	r.addExecution(caller)
	if condition() {
		r.markAsExecuted(caller)
	}
}

func (r *Routine) WaitUntilOrFor(condition func() bool, duration time.Duration) {
	if !r.started {
		return
	}

	caller := r.caller()
	if r.isExecuted(caller) {
		return
	}
	if !r.isPrevExecuted(caller) {
		return
	}

	timer := r.executionTimer(caller, duration)
	select {
	case <-timer:
		r.markAsExecuted(caller)
	default:
		if condition() {
			r.markAsExecuted(caller)
		}
	}
}
