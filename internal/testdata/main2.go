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
			fmt.Println("Start")
		})

		for i := 0; i < 3; i++ {
			r.Do(func() {
				fmt.Println("Iteration", i+1, "at", time.Now().Format(time.TimeOnly))
			})
			r.WaitFor(time.Second)
		}

		r.Do(func() {
			fmt.Println("End")
		})

		r.End()
	}
	fmt.Println("Loops", loops)
}
