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
			fmt.Println("Before", time.Now().Format(time.TimeOnly))
		})

		r.WaitUntil(func() bool {
			return loops == 10
		})

		r.Do(func() {
			fmt.Println("Current iteration", loops)
		})

		r.WaitFor(time.Second)

		r.Do(func() {
			fmt.Println("After", time.Now().Format(time.TimeOnly))
		})

		r.WaitUntilOrFor(func() bool {
			return true
		}, time.Second)

		r.Do(func() {
			fmt.Println("Done", time.Now().Format(time.TimeOnly))
		})

		r.End()
	}
	fmt.Println("Loops", loops)

	// fmt.Println()
	// fmt.Println("Execution sequence:")
	// fmt.Println(r.ExecutionSequenceRaw())
	// fmt.Println(strings.Join(r.ExecutionSequence(), "\n"))
}
