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

		const n = 2
		const t = time.Second / 16
		r.Do(func() {
			fmt.Println("=")
		})
		r.WaitFor(t)
		r.Func(func() {
			r.Repeat(n, func() {
				r.Do(func() {
					fmt.Println("+")
				})
				r.WaitFor(t)
			})
		})
		r.WaitFor(t)
		r.Do(func() {
			fmt.Println("=")
		})
		r.Func(func() {
			r.Repeat(n, func() {
				r.Do(func() {
					fmt.Println("-")
				})
				r.WaitFor(t)
			})
		})
		r.WaitFor(t)
		r.Do(func() {
			fmt.Println("=")
		})

		r.End()
	}
	fmt.Println("Loops", loops)
}
