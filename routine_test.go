package routines_test

import (
	"testing"
	"time"

	"github.com/mymmrac/routines"
	"github.com/mymmrac/routines/internal/test"
)

const maxLoop = 400
const waitTime = time.Millisecond

func TestRoutine_StartEnd(t *testing.T) {
	r := routines.NewRoutine()

	test.False(t, r.Started())
	test.False(t, r.Completed())

	r.End()

	test.False(t, r.Started())
	test.False(t, r.Completed())

	r.Start()

	test.True(t, r.Started())
	test.False(t, r.Completed())

	r.End()

	test.False(t, r.Started())
	test.True(t, r.Completed())

	r.Reset()

	test.False(t, r.Started())
	test.False(t, r.Completed())
}

func TestStartRoutine(t *testing.T) {
	r := routines.StartRoutine()

	test.True(t, r.Started())
	test.False(t, r.Completed())
}

func TestRoutine_Do(t *testing.T) {
	r := routines.NewRoutine()

	t.Run("not-started", func(t *testing.T) {
		r.Do(func() {
			t.FailNow()
		})
		r.Func(func() {
			t.FailNow()
		})
		r.Loop(0, 1, func(_ int) {
			t.FailNow()
		})
		r.Repeat(1, func() {
			t.FailNow()
		})
	})

	r.Start()
	test.True(t, r.Started())

	t.Run("do", func(t *testing.T) {
		e1 := false
		r.Do(func() {
			e1 = true
		})
		test.True(t, e1)
	})

	t.Run("do-many", func(t *testing.T) {
		e1 := false
		r.Do(func() {
			e1 = true
		})
		test.True(t, e1)

		e2 := false
		r.Do(func() {
			e2 = true
		})
		test.True(t, e2)

		e3 := false
		r.Do(func() {
			e3 = true
		})
		test.True(t, e3)
	})

	t.Run("do-once", func(t *testing.T) {
		e1 := 0
		for i := 0; i < 2; i++ {
			r.Do(func() {
				e1++
			})
		}
		test.Equal(t, e1, 1)
	})

	test.True(t, r.Started())
	test.False(t, r.Completed())

	r.End()

	test.False(t, r.Started())
	test.True(t, r.Completed())
}

func TestDoAndWait(t *testing.T) {
	r := routines.NewRoutine()

	e1 := false
	e2 := false
	e3 := false
	e4 := false

	loops := 0
	var loopsOnTrue int
	for !r.Completed() && loops < maxLoop {
		loops++

		r.Start()

		r.Do(func() {
			test.False(t, e2 || e3 || e4)
			test.Equal(t, loops, 1)
			e1 = true
		})

		r.WaitUntil(func() bool {
			return loops == 10
		})

		r.Do(func() {
			test.True(t, e1)
			test.False(t, e3 || e4)
			test.Equal(t, loops, 10)
			e2 = true
		})

		r.WaitFor(waitTime)

		r.Do(func() {
			test.True(t, e1 && e2)
			test.False(t, e4)
			test.True(t, loops > 10)
			e3 = true

			loopsOnTrue = loops
		})

		r.WaitUntilOrTimeout(func() bool {
			return true
		}, time.Hour)

		r.Do(func() {
			test.True(t, e1 && e2 && e3)
			test.Equal(t, loops, loopsOnTrue)
			e4 = true
		})

		r.End()
	}

	test.True(t, loops < maxLoop)
	test.True(t, e1 && e2 && e3 && e4)
}

func TestRoutine_Loop(t *testing.T) {
	r := routines.StartRoutine()

	e1 := false
	e2 := 0
	e3 := false

	loops := 0
	for !r.Completed() && loops < maxLoop {
		loops++

		r.Do(func() {
			test.False(t, e1)
			e1 = true
		})

		r.Loop(0, 3, func(i int) {
			r.Do(func() {
				test.True(t, e1)
				e2++
			})
			r.WaitFor(waitTime)
		})

		r.Do(func() {
			test.True(t, e1)
			test.Equal(t, e2, 3)
			e3 = true
		})

		r.End()
	}

	test.True(t, loops < maxLoop)
	test.True(t, e1 && e2 == 3 && e3)
}

func TestNestedLoop(t *testing.T) {
	r := routines.StartRoutine()

	e1 := false
	e2 := 0
	var e3 []int
	e4 := false

	loops := 0
	for !r.Completed() && loops < maxLoop {
		loops++

		r.Do(func() {
			test.False(t, e1)
			e1 = true
		})
		r.WaitFor(waitTime)

		r.Loop(0, 3, func(i int) {
			r.Do(func() {
				test.True(t, e1)
				e2++

				e3 = append(e3, 0)
			})
			r.WaitFor(waitTime)

			r.Loop(0, 2, func(j int) {
				r.Do(func() {
					e3[i]++
				})
				r.WaitFor(waitTime)
			})
		})

		r.Do(func() {
			test.True(t, e1)
			test.Equal(t, e2, 3)
			test.EqualEl(t, e3, []int{2, 2, 2})
			e4 = true
		})

		r.End()
	}

	test.True(t, loops < maxLoop)
	test.True(t, e1 && e2 == 3 && len(e3) == 3 && e4)
}

func TestNestedDoAndWait(t *testing.T) {
	r := routines.NewRoutine()

	e1 := false
	e2 := false
	e3 := false
	e4 := false
	e5 := 0

	doStuff := func(r *routines.Routine) {
		r.Do(func() {
			e5++
		})
		// This can't include any waiters, loop will freeze
		// r.Func should be used
	}

	loops := 0
	for !r.Completed() && loops < maxLoop {
		loops++

		r.Start()

		r.Do(func() {
			test.False(t, e1 || e2 || e3 || e4)
			e1 = true
		})

		r.WaitFor(waitTime)

		r.Do(func() {
			test.False(t, e2 || e3 || e4)
			test.True(t, e1)
			e2 = true
			test.Equal(t, e5, 0)
			doStuff(r)
			test.Equal(t, e5, 1)
		})

		r.WaitFor(waitTime)

		r.Do(func() {
			test.False(t, e3 || e4)
			test.True(t, e1 && e2)
			e3 = true
		})

		r.WaitFor(waitTime)

		r.Do(func() {
			test.False(t, e4)
			test.True(t, e1 && e2 && e3)
			e4 = true
			test.Equal(t, e5, 1)
			doStuff(r)
			test.Equal(t, e5, 2)
		})

		r.End()
	}

	test.True(t, loops < maxLoop)
	test.True(t, e1 && e2 && e3 && e4 && e5 == 2)
}

func TestRoutine_Restart(t *testing.T) {
	r := routines.StartRoutine()
	test.True(t, r.Started())

	e1 := false
	e2 := false

	loops := 0
	for !r.Completed() && loops < maxLoop {
		loops++

		r.Do(func() {
			e1 = true
		})

		r.WaitFor(waitTime)
		r.End()
	}
	test.True(t, loops < maxLoop)
	test.True(t, e1)
	test.True(t, r.Completed())

	r.Restart()
	test.True(t, r.Started())
	test.False(t, r.Completed())

	loops = 0
	for !r.Completed() && loops < maxLoop {
		loops++

		r.Do(func() {
			e2 = true
		})

		r.WaitFor(waitTime)
		r.End()
	}
	test.True(t, loops < maxLoop)
	test.True(t, e2)
	test.True(t, r.Completed())
}

func TestRoutine_WaitForDone(t *testing.T) {
	r := routines.NewRoutine()

	e1 := false
	e2 := false
	e3 := false

	wait := make(chan struct{})
	neverDone := make(chan struct{})

	loops := 0
	for !r.Completed() && loops < maxLoop {
		loops++

		r.Start()

		r.Do(func() {
			close(wait)
			test.False(t, e1 || e2 || e3)
			e1 = true
		})

		r.WaitForDone(wait)

		r.Do(func() {
			test.False(t, e2 || e3)
			test.True(t, e1)
			e2 = true
		})

		r.WaitForDoneOrTimeout(neverDone, waitTime)

		r.Do(func() {
			test.False(t, e3)
			test.True(t, e1 && e2)
			e3 = true
		})

		r.End()
	}
	test.True(t, loops < maxLoop)
	test.True(t, e1 && e2 && e3)
}

func TestRoutine_Func(t *testing.T) {
	r := routines.NewRoutine()

	e1 := false
	e2 := false
	e3 := 0
	e4 := 0

	do := func(r *routines.Routine, i int) {
		r.Do(func() {
			switch i {
			case 1:
				test.False(t, e2)
				test.True(t, e1)
				test.Equal(t, e3, 0)
				test.Equal(t, e4, 0)
			case 2:
				test.False(t, e2)
				test.True(t, e1)
				test.Equal(t, e3, 1)
				test.Equal(t, e4, 1)
			}
			e3++
		})

		r.WaitFor(waitTime)

		r.Do(func() {
			switch i {
			case 1:
				test.False(t, e2)
				test.True(t, e1)
				test.Equal(t, e3, 1)
				test.Equal(t, e4, 0)
			case 2:
				test.False(t, e2)
				test.True(t, e1)
				test.Equal(t, e3, 2)
				test.Equal(t, e4, 1)
			}
			e4++
		})
	}

	loops := 0
	for !r.Completed() && loops < maxLoop {
		loops++

		r.Start()

		r.Do(func() {
			test.False(t, e1 || e2)
			e1 = true
		})

		r.Func(func() {
			do(r, 1)
		})

		r.WaitFor(waitTime)

		r.Func(func() {
			do(r, 2)
		})

		r.Do(func() {
			test.False(t, e2)
			test.True(t, e1)
			e2 = true
		})

		r.End()
	}
	test.True(t, loops < maxLoop)
	test.True(t, e1 && e2 && e3 == 2 && e4 == 2)
}

func TestRoutine_Repeat(t *testing.T) {
	r := routines.NewRoutine()

	e1 := false
	e2 := false
	e3 := 0
	e4 := 0
	e5 := 0

	doRepeated := func(r *routines.Routine) {
		r.Do(func() {
			e3++
		})

		r.WaitFor(waitTime)

		r.Repeat(3, func() {
			r.Do(func() {
				e4++
			})
			r.WaitFor(waitTime)
		})

		r.WaitFor(waitTime)

		r.Do(func() {
			e5++
		})
	}

	loops := 0
	for !r.Completed() && loops < maxLoop {
		loops++

		r.Start()

		r.Do(func() {
			test.False(t, e1 || e2)
			e1 = true
		})

		r.Repeat(2, func() {
			doRepeated(r)
		})

		r.Do(func() {
			test.False(t, e2)
			test.True(t, e1)
			e2 = true
		})

		r.End()
	}
	test.True(t, loops < maxLoop)
	test.True(t, e1 && e2 && e3 == 2 && e4 == 6 && e5 == 2)
}
