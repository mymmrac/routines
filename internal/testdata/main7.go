package main

import (
	"fmt"
	"time"

	"github.com/mymmrac/routines"
)

func main() {
	r := routines.NewRoutine()

	loops := 0
	for !r.Completed() {
		loops++

		r.Start()

		r.Do(func() {
			fmt.Println("Before")
		})

		r.Func(func() {
			do(r, 1)
		})

		r.Func(func() {
			do(r, 2)
		})

		r.Do(func() {
			fmt.Println("After")
		})

		r.End()
	}
	fmt.Println("Loops", loops)
}

func do(r *routines.Routine, i int) {
	r.Do(func() {
		fmt.Println("Before nested", i)
	})

	r.WaitFor(time.Second / 2)

	r.Do(func() {
		fmt.Println("After nested", i)
	})
}
