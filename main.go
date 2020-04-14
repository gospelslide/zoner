package main

import (
	"fmt"
	"github.com/gospelslide/zoner/geoindex"
	"github.com/gospelslide/zoner/workerpool"
	_"bufio"
	_"os"
)

func writeLocations(done chan bool) {
	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	workerpool.LocationWriteQueue <- scanner.Text()
	// }
	for i := 0; i < 20; i++ {
		workerpool.LocationWriteQueue <- geoindex.Location{Lat: 19.098874, Long: 72.908818}
	}
	done <- true
}

func readIndexedLocations(done chan bool) {
	for result := range workerpool.IndexedLocationReadQueue {
		fmt.Println(result)
	}
	done <- true
}

func main() {
	done := make(chan bool)
	go writeLocations(done)
	go readIndexedLocations(done)
	workerpool.CreateLocationIndexWorkerPool(20)
	<- done
	// loc := geoindex.Location{Lat: 19.098874, Long: 72.908818}
	// fmt.Println(geoindex.LocationToGeoIndex(loc, 3).Index)
}