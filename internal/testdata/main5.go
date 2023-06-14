package main

import (
	"fmt"
	"time"

	"github.com/mymmrac/routines"
)

func main() {
	r := routines.StartRoutine()

	loops := 0
	for !r.Completed() {
		loops++

		r.Do(func() {
			fmt.Println("Stuff")
		})

		r.WaitFor(time.Millisecond)
		r.End()
	}
	fmt.Println("Loops", loops)

	r.Restart()
	for !r.Completed() {
		loops++

		r.Do(func() {
			fmt.Println("Other stuff")
		})

		r.WaitFor(time.Millisecond)
		r.End()
	}
	fmt.Println("Loops", loops)
}
