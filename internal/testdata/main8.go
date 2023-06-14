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

		r.Repeat(2, func() {
			doRepeated(r)
		})

		r.Do(func() {
			fmt.Println("After")
		})

		r.End()
	}
	fmt.Println("Loops", loops)
}

func doRepeated(r *routines.Routine) {
	r.Do(func() {
		fmt.Println("Before nested")
	})

	r.WaitFor(time.Second / 2)

	r.Repeat(5, func() {
		r.Do(func() {
			fmt.Print(".")
		})
		r.WaitFor(time.Second / 5)
	})
	r.Do(func() {
		fmt.Println()
	})

	r.Do(func() {
		fmt.Println("After nested")
	})
}
