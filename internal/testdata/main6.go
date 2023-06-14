package main

import (
	"fmt"
	"time"

	"github.com/mymmrac/routines"
)

func main() {
	r := routines.NewRoutine()

	var input string
	wait := make(chan struct{})
	done := make(chan struct{})
	go func() {
		<-wait
		fmt.Print("> ")
		_, _ = fmt.Scanln(&input)
		close(done)
	}()

	neverDone := make(chan struct{})

	loops := 0
	for !r.Completed() {
		loops++

		r.Start()

		r.Do(func() {
			fmt.Println("Start")
			close(wait)
		})

		r.WaitForDone(done)

		r.Do(func() {
			fmt.Println("Input:", input)
		})

		r.WaitForDoneOrTimeout(neverDone, time.Second/2)

		r.Do(func() {
			fmt.Println("Done")
		})

		r.End()
	}
	fmt.Println("Loops", loops)
}
