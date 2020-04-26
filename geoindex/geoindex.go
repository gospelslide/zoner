package geoindex

import (
	"fmt"
	"math"
)

type Location struct {
	Lat, Long float64
	Index string
}

type Quadrant struct {
	Lat0, Long0, LatBottom, LatTop, LongLeft, LongRight float64
	Index                                               string
}

var baseIterations = 3 // represents times to divide the start point of each map i.e 3 -> 4^3 -> 1/64th size of each quadrant
var gridBlocks = 1 << baseIterations
var locationQueue = make(chan Location, 100)
var locationIndexQueue = make(chan Location, 100)

var fullMap = Quadrant{
	Lat0:      0.00,
	Long0:     0.00,
	LongLeft:  -180.00,
	LongRight: 180.00,
	LatTop:    90.00,
	LatBottom: -90.00,
	Index:     "",
}

func latitudeWithTopAsAxis(latitude float64, q Quadrant) float64 {
	return math.Abs(latitude - q.LatBottom)
}

func longitudeWithLeftAsAxis(longitude float64, q Quadrant) float64 {
	return math.Abs(longitude - q.LongLeft)
}

func checkIfMapSpreadAcrossAxes(x, y float64) bool {
	return (x < 0 && y > 0)
}

func mapLongToNumber(longitude float64, q Quadrant) int {
	longitude = longitudeWithLeftAsAxis(longitude, q)
	sumLongitude := math.Abs(q.LongLeft - q.LongRight)
	longitudinalBlock := sumLongitude / float64(gridBlocks)
	// fmt.Println("long-", longitude, sumLongitude, longitudinalBlock)
	return int(math.Floor(longitude / longitudinalBlock))
}

func mapLatToAlpha(latitude float64, q Quadrant) string {
	latitude = latitudeWithTopAsAxis(latitude, q)
	sumLatitude := math.Abs(q.LatBottom - q.LatTop)
	latitudinalBlock := sumLatitude / float64(gridBlocks)
	// fmt.Println("lat-", latitude, sumLatitude, latitudinalBlock)
	return string(int(math.Floor(latitude/latitudinalBlock)) + 97)
}

func calculateIndexForQuadrant(q, prevGranularity Quadrant) string {
	latitudeAlpha := mapLatToAlpha(q.Lat0, prevGranularity)
	longitudeNo := mapLongToNumber(q.Long0, prevGranularity)
	return fmt.Sprintf("%s%d", latitudeAlpha, longitudeNo)
}

func divideQuadrants(q, prevGranularity Quadrant, point Location, iter int) Quadrant {
	if iter == 0 {
		// fmt.Println("break-", prevGranularity, q)
		q.Index += calculateIndexForQuadrant(q, prevGranularity)
		// fmt.Println("index-", q.Index)
		return q
	}
	newq := q
	if point.Lat > q.Lat0 {
		newq.Lat0 = (q.Lat0 + q.LatTop) / 2
		newq.LatBottom = q.Lat0
	} else {
		newq.Lat0 = (q.Lat0 + q.LatBottom) / 2
		newq.LatTop = q.Lat0
	}
	if point.Long > q.Long0 {
		newq.Long0 = (q.Long0 + q.LongRight) / 2
		newq.LongLeft = q.Long0
	} else {
		newq.Long0 = (q.Long0 + q.LongLeft) / 2
		newq.LongRight = q.Long0
	}
	iter--
	// fmt.Println("curr-", q)
	return divideQuadrants(newq, prevGranularity, point, iter)
}

/*LocationToGeoIndex function to 
@params - Location struct, int granularity of index
@return - Quadrant struct with axes, lat long boundaries and index
*/
func LocationToGeoIndex(point Location, granularity int) Location {
	q := fullMap
	for i := 0; i < granularity; i++ {
		q = divideQuadrants(q, q, point, baseIterations)
	}
	point.Index = q.Index
	return point
}

/*
Tests
19.098874, 72.908818 - e5g4g7
34.857166, 76.959279 - f5e5d5
*/
