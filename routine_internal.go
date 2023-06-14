package routines

import (
	"fmt"
	"runtime"
	"time"
)

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
	return encodeCaller(r.executionStack), func() {
		r.executionStack = r.executionStack[:len(r.executionStack)-1]
	}
}
