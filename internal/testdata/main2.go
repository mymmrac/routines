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
			fmt.Println("Start")
		})

		r.For(0, 3, func(i int) {
			r.Do(func() {
				fmt.Println("Iteration", i+1, "at", time.Now().Format(time.TimeOnly))
			})
			r.WaitFor(time.Second)
		})

		r.Do(func() {
			fmt.Println("End")
		})

		r.End()
	}
	fmt.Println("Loops", loops)
}
