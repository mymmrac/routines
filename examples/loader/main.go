package main

import (
	"fmt"
	"time"

	"github.com/mymmrac/routines"
)

func main() {
	r := routines.StartRoutine()
	for !r.Completed() {
		r.Do(func() {
			fmt.Println("Hello Routines!")
			fmt.Print("Loading")
		})
		r.Repeat(3, func() {
			r.WaitFor(time.Second / 2)
			r.Do(func() {
				fmt.Print(".")
			})
		})
		r.Do(func() {
			fmt.Println()
			fmt.Println("Done!")
		})
		r.End()
	}
}
