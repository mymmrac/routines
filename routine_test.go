package routines_test

import (
	"testing"
	"time"

	"github.com/mymmrac/routines"
	"github.com/mymmrac/routines/internal/test"
)

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
		executed := false
		r.Do(func() {
			executed = true
		})
		test.True(t, executed)
	})

	t.Run("do-many", func(t *testing.T) {
		executed1 := false
		r.Do(func() {
			executed1 = true
		})
		test.True(t, executed1)

		executed2 := false
		r.Do(func() {
			executed2 = true
		})
		test.True(t, executed2)

		executed3 := false
		r.Do(func() {
			executed3 = true
		})
		test.True(t, executed3)
	})

	t.Run("do-once", func(t *testing.T) {
		executed := 0
		for i := 0; i < 2; i++ {
			r.Do(func() {
				executed++
			})
		}
		test.Equal(t, executed, 1)
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
	for !r.Completed() && loops < 10000 {
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

		r.WaitFor(time.Millisecond)

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

	test.True(t, e1 && e2 && e3 && e4)
}
