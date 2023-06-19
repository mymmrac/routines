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

		r.WaitFor(time.Millisecond)

		r.Do(func() {
			doStuff(r)
		})

		r.WaitFor(time.Millisecond)

		r.Do(func() {
			fmt.Println("After")
		})

		r.WaitFor(time.Millisecond)

		r.Do(func() {
			doStuff(r)
		})

		r.WaitFor(time.Millisecond)
		r.End()
	}
	fmt.Println("Loops", loops)
}

func doStuff(r *routines.Routine) {
	r.Do(func() {
		fmt.Println("The Stuff")
	})
}
