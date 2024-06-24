package main

import (
	"fmt"
	"log"
	"nyct-arrival-times"
)

func main() {
	nyctStopId := 354
	nyctDirection := "N"
	trainsToReturn := 2
	response, err := nyct.CheckArrivalTimes(nyctStopId, nyctDirection, trainsToReturn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}
