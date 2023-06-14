package main

import (
	"fmt"
	"time"

	"github.com/mymmrac/routines"
)

const timeFormat = "15:04:05.999"

func main() {
	r := routines.StartRoutine()

	loops := 0
	for !r.Completed() {
		loops++

		r.Do(func() {
			fmt.Println("Start")
		})

		r.Loop(0, 3, func(i int) {
			r.Do(func() {
				fmt.Println("Iteration", i+1, "at", time.Now().Format(timeFormat))
			})
			r.WaitFor(time.Second / 2)

			r.Loop(0, 2, func(j int) {
				r.Do(func() {
					fmt.Println("Iteration", i+1, ":", j+1, "at", time.Now().Format(timeFormat))
				})
				r.WaitFor(time.Second / 2)
			})
		})

		r.Do(func() {
			fmt.Println("End")
		})

		r.End()
	}
	fmt.Println("Loops", loops)
}
