package routines

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

func (r *Routine) Restart() {
	r.Reset()
	r.Start()
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

func (r *Routine) Func(action func()) {
	if !r.started {
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

	action()

	if r.isPrevExecuted(caller) {
		r.addExecution(caller)
		r.markAsExecuted(caller)
	}
}

func (r *Routine) Loop(start, end int, action func(i int)) {
	if !r.started {
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

	if r.isPrevExecuted(caller) {
		r.addExecution(caller)
		r.markAsExecuted(caller)
	}
}

func (r *Routine) Repeat(n int, action func()) {
	if !r.started {
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

	for i := 0; i < n; i++ {
		_, popIndex := r.pushToStack(uintptr(i))
		action()
		popIndex()
	}

	if r.isPrevExecuted(caller) {
		r.addExecution(caller)
		r.markAsExecuted(caller)
	}
}
