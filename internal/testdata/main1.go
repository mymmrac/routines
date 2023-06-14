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

		r.WaitFor(time.Second * 2)

		r.Do(func() {
			fmt.Println("After", time.Now().Format(time.TimeOnly))
		})

		r.WaitUntilOrTimeout(func() bool {
			return true
		}, time.Second)

		r.Do(func() {
			fmt.Println("Done", time.Now().Format(time.TimeOnly))
		})

		r.End()
	}
	fmt.Println("Loops", loops)
}
