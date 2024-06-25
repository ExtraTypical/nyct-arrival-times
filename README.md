# NYCT (MTA) Subway Arrival Times

A Go package for the express purpose of checking arrival times of trains on the NYCT subway.

I was originally searching for a package like this in Go or JavaScript but couldn't find anything that was up to date with the new MTA Developer changes. This is my first package so I'm sure there are inefficiencies, but I'm sure someone will find it useful.

This package also means that you won't have to work with the gtfs-realtime.proto or nyct-subway.proto files.

## Installation

In your Go project, start with the standard init for your Go mod manager.

```
go mod init your-project
```

Once you've done that, run the command:

```
go mod tidy
```

Then run the command:

```
go get github.com/ExtraTypical/nyct-arrival-times@v1.0.4
```

This will download the package and allow you to use it in your code.

## Usage

There's one public function with this package called ```CheckArrivalTimes```.

```go
package main

func main() {

    nyctStopId := "235"
    nyctDirection := "N" /* Optional: Can also just be an empty string */
    trainsToReturn := 4

    response, err := nyct.CheckArrivalTimes(nyctStopId, nyctDirection, trainsToReturn)
    if err != nil {
        /* Your error handling here */
    }

    fmt.Println(response)
}
```

Get your nyctStopId by heading to [this website](https://data.ny.gov/Transportation/MTA-Subway-Stations/39hk-dx4f/data_preview) and searching for the GTFS ID of your stop.