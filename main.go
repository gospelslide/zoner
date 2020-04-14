package main

import (
	"fmt"
	"github.com/gospelslide/zoner/geoindex"
	"github.com/gospelslide/zoner/workerpool"
	_"bufio"
	_"os"
	"time"
)

func writeLocations() {
	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	workerpool.LocationWriteQueue <- scanner.Text()
	// }
	for i := 0; ; i++ {
		workerpool.LocationWriteQueue <- geoindex.Location{Lat: 19.098874, Long: 72.908818}
	}
}

func readIndexedLocations(done chan bool) {
	counter := 0
	startTime := time.Now()
	for result := range workerpool.IndexedLocationReadQueue {
		counter++
		fmt.Println(result)
		if counter == 50 {
			timeTaken := time.Now().Sub(startTime)
			fmt.Println("Total time taken for 50 logs - ", timeTaken.Microseconds(), "ms")
			counter = 0
			startTime = time.Now()
		}
	}
	done <- true
}

func main() {
	done := make(chan bool)
	go writeLocations()
	go readIndexedLocations(done)
	workerpool.CreateLocationIndexWorkerPool(20)
	<- done
	// loc := geoindex.Location{Lat: 19.098874, Long: 72.908818}
	// fmt.Println(geoindex.LocationToGeoIndex(loc, 3).Index)
}