package main

import (
	"fmt"
)

func main() {
	jobs := []job{
		job(func(in, out chan interface{}) {
			out <- 0
			out <- 1
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("Error!")
			}
			fmt.Println("Result is: ", data)
		}),
	}

	ExecutePipeline(jobs...)
}
