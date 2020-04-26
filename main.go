package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/gospelslide/zoner/geoindex"
	_ "github.com/gospelslide/zoner/workerpool"
	"github.com/Shopify/sarama"
	"log"
	"errors"
)

var defaultGranularity int64 = 3

type Server struct {
	LocationDataProducer sarama.AsyncProducer
}

type LocationData struct {
	Lat float64     `json:"latitude"`
	Long float64    `json:"longitude"`
	GeoIndex string    `json:"geoindex"`
	LocType string   `json:"location_type"`
	GeoIndexGranularity int8 `json:"geoindex_granularity"`
}

func (s *Server) Close() error {
	err := s.LocationDataProducer.Close()
	if err != nil {
		log.Println("Failed to shut down producer cleanly", err)
	}
	return nil
}

func (s *Server) Run(addr string) error {
	httpServer := http.Server{
		Addr: addr,
		Handler: s.Handler(),
	}
	log.Printf("Listening for requests on %s....\n", addr)
	return httpServer.ListenAndServe()
}

func (s *Server) Handler() http.Handler {
	return s.locationDataHandler()
}

func (s *Server) locationDataHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message, err := indexAndMarshalData(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid data")
			return
		}
		s.LocationDataProducer.Input() <- &sarama.ProducerMessage{
			Topic: "important",
			Key: sarama.StringEncoder("1"),
			Value: sarama.StringEncoder(message),
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ok")
	})
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

func newLocationDataProducer(brokerlist []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerlist, config)
	if err != nil {
		log.Fatalln("Failed to initiate location data producer:", err)
	}

	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write location data:", err)
		}
	}()
	return producer
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
	brokers := []string{"0.0.0.0:9092"}
	addr := ":8080"
	server := &Server{
		LocationDataProducer: newLocationDataProducer(brokers),
	}
	defer func() {
		err := server.Close()
		if err != nil {
			log.Println("Failed to cleanup and close server", err)
		}
	}()
	log.Fatal(server.Run(addr))
}