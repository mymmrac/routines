package routines

import "time"

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
