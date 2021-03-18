package main

import "fmt"

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

//	jobs := []job{
//		job(func(in, out chan interface{}) {
//			fmt.Println("-1")
//		}),
//		job(func(in, out chan interface{}) {
//			fmt.Println("-2")
//		}),
//		job(func(in, out chan interface{}) {
//			fmt.Println("-3")
//		}),
//		job(func(in, out chan interface{}) {
//			fmt.Println("-4")
//		}),
//		job(func(in, out chan interface{}) {
//			fmt.Println("-5")
//		}),
//	}

	ExecutePipeline(jobs...)
}
