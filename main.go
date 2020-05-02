package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gospelslide/zoner/geoindex"
	_ "github.com/gospelslide/zoner/workerpool"
	"log"
	"errors"
	"strconv"
)

var defaultGranularity int64 = 3

type LocationData struct {
	Lat float64     `json:"latitude"`
	Long float64    `json:"longitude"`
	GeoIndex string    `json:"geoindex"`
	LocType string   `json:"location_type"`
	GeoIndexGranularity int8 `json:"geoindex_granularity"`
}

func locationDataHandler(w http.ResponseWriter, r *http.Request) {
	message, err := indexAndMarshalData(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid data")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(message)
}

func indexAndMarshalData(r *http.Request) ([]byte, error) {
	var lat, long float64
	var err error
	var granularity int64
	lat, err = strconv.ParseFloat(readQueryParam("lat", r), 64)
	long, err = strconv.ParseFloat(readQueryParam("long", r), 64)
	locType := readQueryParam("locType", r)
	granularity, intErr := strconv.ParseInt(readQueryParam("granularity", r), 10, 32)
	if intErr != nil {
		granularity = defaultGranularity
	}
	if err != nil || locType == "" {
		return nil, errors.New("Data format invalid")
	}
	locationData := LocationData{
		Lat: lat,
		Long: long,
		LocType: locType,
		GeoIndexGranularity: int8(granularity),
	}
	locationData = geoIndexLocation(locationData)
	return json.Marshal(locationData)
}

func readQueryParam(param string, r *http.Request) string {
	data := r.URL.Query().Get(param)
	if (data != "") { 
		return string(data[0]) 
	}
	return ""
}

func geoIndexLocation(data LocationData) LocationData {
	point := geoindex.Location{
		Lat: data.Lat,
		Long: data.Long,
	}
	point = geoindex.LocationToGeoIndex(point, int(data.GeoIndexGranularity))
	data.GeoIndex = point.Index
	return data
}

// func writeLocations(w http.ResponseWriter, req *http.Request) {
// 	lat, err := strconv.ParseFloat(req.URL.Query()["lat"][0], 64)
// 	long, err := strconv.ParseFloat(req.URL.Query()["long"][0], 64)
// 	if err == nil { workerpool.LocationWriteQueue <- geoindex.Location{Lat: lat, Long: long} }
// 	fmt.Fprint(w, "Ok")
// }

// func readIndexedLocations(done chan bool) {
// 	counter := 0
// 	startTime := time.Now()
// 	for result := range workerpool.IndexedLocationReadQueue {
// 		counter++
// 		fmt.Println(result)
// 		if counter == 50 {
// 			timeTaken := time.Now().Sub(startTime)
// 			fmt.Println("Total time taken for 50 logs - ", timeTaken.Microseconds(), "ms")
// 			counter = 0
// 			startTime = time.Now()
// 		}
// 	}
// 	done <- true
// }

// func initialise() {
// 	done := make(chan bool)
// 	go readIndexedLocations(done)
// 	workerpool.CreateLocationIndexWorkerPool(20)
// 	<- done
// }

func main() {
	http.HandleFunc("/", locationDataHandler)
	port := ":8080"
	fmt.Printf("Initialized server at %s ...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}