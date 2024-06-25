package main

import (
	"fmt"
	"log"

	"github.com/ExtraTypical/nyct-arrival-times"
)

func main() {
	nyctStopId := "235"
	nyctDirection := "N"
	trainsToReturn := 4
	response, err := nyct.CheckArrivalTimes(nyctStopId, nyctDirection, trainsToReturn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}
