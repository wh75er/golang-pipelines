package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	MultiHashNumber = 6
)

type multiHashValue struct {
	th   int
	hash string
}

type multiHashSlice []multiHashValue

func (s multiHashSlice) Len() int {
	return len(s)
}

func (s multiHashSlice) Less(i, j int) bool {
	return s[i].th < s[j].th
}

func (s multiHashSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func getMd5Crc32(m *sync.Mutex, data string, c chan<- string) {
	m.Lock()
	md5 := DataSignerMd5(data)
	m.Unlock()
	fmt.Println(data, " SingleHash md5(data) ", md5)

	md5Crc32 := DataSignerCrc32(md5)
	fmt.Println(data, " SingleHash crc32(md5(data)) ", md5Crc32)

	c <- md5Crc32
}

func getSingleHash(m *sync.Mutex, inItem interface{}, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	rawData := inItem

	data := strconv.Itoa(rawData.(int))
	fmt.Println(data, " SingleHash data ", data)

	c := make(chan string)
	go getMd5Crc32(m, data, c)

	crc32 := DataSignerCrc32(data)
	fmt.Println(data, " SingleHash crc32(data) ", crc32)

	result := crc32 + "~" + <-c
	fmt.Println(data, " SingleHash result ", result)

	out <- result
}

func SingleHash(in, out chan interface{}) {
	m := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for inItem := range in {
		wg.Add(1)
		go getSingleHash(m, inItem, out, wg)
	}

	wg.Wait()
}

func concatMultiHashes(r <-chan multiHashValue) string {
	var results multiHashSlice

	for v := range r {
		results = append(results, v)
	}

	sort.Sort(results)

	var result string

	for _, v := range results {
		result += v.hash
	}

	return result
}

func getMultiHash(r chan multiHashValue, th int, data string, wgTh *sync.WaitGroup) {
	defer wgTh.Done()
	step := DataSignerCrc32(strconv.Itoa(th) + data)

	fmt.Println(data, " MultiHash: crc32(th+step1) ", th, step)

	r <- multiHashValue{th, step}
}

func getMultiHashes(inItem interface{}, out chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	rawData := inItem

	data := rawData.(string)

	r := make(chan multiHashValue, MultiHashNumber)
	wgTh := &sync.WaitGroup{}

	for th := 0; th < MultiHashNumber; th++ {
		wgTh.Add(1)
		go getMultiHash(r, th, data, wgTh)
	}

	wgTh.Wait()
	close(r)

	result := concatMultiHashes(r)

	fmt.Println(data, " MultiHash: result ", result)

	out <- result
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for inItem := range in {
		wg.Add(1)
		go getMultiHashes(inItem, out, wg)
	}

	wg.Wait()
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
		go func(j job, in, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()

			j(in, out)

			close(out)
		}(j, in, out, &wg)

		in = out
	}

	wg.Wait()
}
