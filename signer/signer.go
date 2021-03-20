package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

func SingleHash(in, out chan interface{})  {
	m := &sync.Mutex{}

	for inItem := range in {
		rawData := inItem

		data := strconv.Itoa(rawData.(int))
		fmt.Println(data, " SingleHash data ", data)

		c := make(chan string)
		go func(m *sync.Mutex, data string, c chan<- string) {
			defer m.Unlock()
			m.Lock()

			md5 := DataSignerMd5(data)
			fmt.Println(data, " SingleHash md5(data) ", md5)

			md5Crc32 := DataSignerCrc32(md5)
			fmt.Println(data, " SingleHash crc32(md5(data)) ", md5Crc32)

			c <- md5Crc32
		}(m, data, c)

		crc32 := DataSignerCrc32(data)
		fmt.Println(data, " SingleHash crc32(data) ", crc32)

		result := crc32 + "~" + <- c
		fmt.Println(data, " SingleHash result ", result)

		out <- result
	}
}

func MultiHash(in, out chan interface{}) {
	for inItem := range in {
		rawData := inItem

		data := rawData.(string)

		results := make([]chan string, 6)

		for th := 0; th < 6; th++ {
			results[th] = make(chan string)

			go func(r []chan string, th int, data string) {
				step := DataSignerCrc32(strconv.Itoa(th) + data)

				fmt.Println(data, " MultiHash: crc32(th+step1) ", th, step)

				r[th] <- step
			}(results, th, data)
		}

		result := <-results[0] + <-results[1] + <-results[2] + <-results[3] + <-results[4] + <-results[5]

		fmt.Println(data, " MultiHash: result ", result)

		out <- result
	}
}

func CombineResults(in, out chan interface{}) {
	results := make([]string, 0)

	for i := range in {
		results = append(results, i.(string))
	}

	sort.Strings(results)

	result := strings.Join(results, "_")

	fmt.Println("CombineResults", result)

	out <- result
}

func ExecutePipeline(jobs ...job) {
	var in chan interface{}

	var wg sync.WaitGroup

	for _, j := range jobs {
		out := make(chan interface{})

		wg.Add(1)
		go func (j job, in, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()

			var inWg sync.WaitGroup

			inWg.Add(1)
			go func(j job, in, out chan interface{}, inWg *sync.WaitGroup) {
				defer inWg.Done()
				j(in, out)
			}(j, in, out, &inWg)

			inWg.Wait()

			close(out)
		}(j, in, out, &wg)

		in = out
	}

	wg.Wait()
}
