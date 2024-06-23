package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testapi/internal/protobuf"

	"google.golang.org/protobuf/proto"
)

func main() {

	CallNYCT()
}

func CallNYCT() {
	client := &http.Client{}

	NYCTUrl := "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs"
	request, err := http.NewRequest("GET", NYCTUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var nyctData protobuf.FeedMessage
	if err := proto.Unmarshal(data, &nyctData); err != nil {
		log.Fatal(err)
	}

	fmt.Println(nyctData.Entity)
}
