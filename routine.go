package routines

import (
	"fmt"
	"runtime"
	"time"
)

type Routine struct {
	started           bool
	completed         bool
	executionStack    []uintptr
	executionSeqIndex map[string]int
	executionSequence []string
	executed          map[string]struct{}
	timers            map[string]<-chan time.Time
	pc                []uintptr
}

func NewRoutine() *Routine {
	return &Routine{
		executionStack:    make([]uintptr, 0),
		executionSeqIndex: make(map[string]int),
		executionSequence: make([]string, 0),
		executed:          make(map[string]struct{}),
		timers:            make(map[string]<-chan time.Time),
		pc:                make([]uintptr, 1),
	}
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

func (r *Routine) Start() {
	if r.started {
		return
	}

	caller, pop := r.pushToStack(r.caller())
	defer pop()

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

	caller, pop := r.pushToStack(r.caller())
	defer pop()

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
	r.executionStack = make([]uintptr, 0)
	r.executionSeqIndex = make(map[string]int)
	r.executionSequence = make([]string, 0)
	r.executed = make(map[string]struct{})
	r.timers = make(map[string]<-chan time.Time)
}

func (r *Routine) Restart() {
	r.Reset()
	r.Start()
}

func (r *Routine) Do(action func()) {
	if !r.started {
		return
	}

	caller, pop := r.pushToStack(r.caller())
	defer pop()

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

func (r *Routine) For(start, end int, action func(i int)) {
	if !r.started {
		return
	}
	if start > end {
		return
	}

	caller, pop := r.pushToStack(r.caller())
	defer pop()

	if r.isExecuted(caller) {
		return
	}
	if !r.isPrevExecutedTo(r.executionSequenceIndex(caller)) {
		return
	}

	for i := start; i < end; i++ {
		_, popIndex := r.pushToStack(uintptr(i))
		action(i)
		popIndex()
	}
}

func (r *Routine) WaitFor(duration time.Duration) {
	if !r.started {
		return
	}

	caller, pop := r.pushToStack(r.caller())
	defer pop()

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

	caller, pop := r.pushToStack(r.caller())
	defer pop()

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

	caller, pop := r.pushToStack(r.caller())
	defer pop()

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

func (r *Routine) caller() uintptr {
	if runtime.Callers(3, r.pc) != 1 {
		panic(fmt.Errorf("failed to get caller"))
	}

	return r.pc[0]
}

func (r *Routine) isExecuted(callers string) (executed bool) {
	_, executed = r.executed[callers]
	return executed
}

func (r *Routine) isPrevExecuted(myCaller string) bool {
	for _, caller := range r.executionSequence {
		if caller == myCaller {
			break
		}

		if !r.isExecuted(caller) {
			return false
		}
	}
	return true
}

func (r *Routine) isPrevExecutedTo(index int) bool {
	for i, caller := range r.executionSequence {
		if i < index {
			break
		}

		if !r.isExecuted(caller) {
			return false
		}
	}
	return true
}

func (r *Routine) executionSequenceIndex(caller string) int {
	i, ok := r.executionSeqIndex[caller]
	if !ok {
		i = len(r.executionSequence)
		r.executionSeqIndex[caller] = i
	}
	return i
}

func (r *Routine) addExecution(caller string) {
	for i := len(r.executionSequence) - 1; i >= 0; i-- {
		if r.executionSequence[i] == caller {
			return
		}
	}

	r.executionSequence = append(r.executionSequence, caller)
}

func (r *Routine) markAsExecuted(caller string) {
	r.executed[caller] = struct{}{}
}

func (r *Routine) executionTimer(caller string, duration time.Duration) <-chan time.Time {
	if timer, found := r.timers[caller]; found {
		return timer
	}

	timer := time.After(duration)
	r.timers[caller] = timer
	r.addExecution(caller)
	return timer
}

func (r *Routine) pushToStack(caller uintptr) (string, func()) {
	r.executionStack = append(r.executionStack, caller)
	return hashCaller(r.executionStack), func() {
		r.executionStack = r.executionStack[:len(r.executionStack)-1]
	}
}
