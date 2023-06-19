package routines

import "time"

type Routine struct {
	started           bool
	completed         bool
	executionStack    []uintptr
	executionSeqIndex map[string]int
	executionSequence []string
	executed          map[string]struct{}
	timers            map[string]<-chan time.Time
	pc                [1]uintptr
}

func NewRoutine() *Routine {
	routine := &Routine{
		pc: [1]uintptr{},
	}
	routine.Reset()
	return routine
}

func (r *Routine) Reset() {
	r.started = false
	r.completed = false
	r.executionStack = make([]uintptr, 0)
	r.executionSeqIndex = make(map[string]int)
	r.executionSequence = make([]string, 0)
	r.executed = make(map[string]struct{})
	r.timers = make(map[string]<-chan time.Time)
}

func StartRoutine() *Routine {
	routine := NewRoutine()
	routine.Start()
	return routine
}

func (r *Routine) Started() bool {
	return r.started
}

func (r *Routine) Completed() bool {
	return r.completed
}
