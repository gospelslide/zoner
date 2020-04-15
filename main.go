package main

import (
	"fmt"
	"github.com/gospelslide/zoner/geoindex"
	"github.com/gospelslide/zoner/workerpool"
	"time"
	"net/http"
	"strconv"
)

func writeLocations(w http.ResponseWriter, req *http.Request) {
	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	workerpool.LocationWriteQueue <- scanner.Text()
	// }
	lat, err := strconv.ParseFloat(req.URL.Query()["lat"][0], 64)
	long, err := strconv.ParseFloat(req.URL.Query()["long"][0], 64)
	if (err == nil) { workerpool.LocationWriteQueue <- geoindex.Location{Lat: lat, Long: long} }
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
	http.HandleFunc("/loc", writeLocations)
	http.ListenAndServe(":8080", nil)
	done := make(chan bool)
	go readIndexedLocations(done)
	workerpool.CreateLocationIndexWorkerPool(20)
	<- done
	// loc := geoindex.Location{Lat: 19.098874, Long: 72.908818}
	// fmt.Println(geoindex.LocationToGeoIndex(loc, 3).Index)
}