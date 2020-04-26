package workerpool

import (
	"sync"
	"github.com/gospelslide/zoner/geoindex"
)

var LocationWriteQueue = make(chan geoindex.Location, 100)
var IndexedLocationReadQueue = make(chan geoindex.Location, 100)

func locationIndexWorker(wg *sync.WaitGroup) {
	for unIndexedLocations := range LocationWriteQueue {
		IndexedLocationReadQueue <- geoindex.LocationToGeoIndex(unIndexedLocations, 3)
	}
}

func CreateLocationIndexWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		go locationIndexWorker(&wg)
		wg.Add(1)
	}
	wg.Wait()
	close(IndexedLocationReadQueue)
}
