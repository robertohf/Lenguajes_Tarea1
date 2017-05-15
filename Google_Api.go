package main

import (
	"github.com/kr/pretty"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"log"
)

func main() {
	direction()
	//geocode()
}

func direction() {
	client, err := maps.NewClient(maps.WithAPIKey("AIzaSyBmelZAhVTODrw_gjtueTuHEs9Aka_z9nM"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DirectionsRequest{

		Origin:      "San Pedro Sula",
		Destination: "Tegucigalpa",
		Mode:        maps.TravelModeDriving,
	}
	routes, _, err := client.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	pretty.Println(routes)
}

func geocode() {
	client, err := maps.NewClient(maps.WithAPIKey("AIzaSyBmelZAhVTODrw_gjtueTuHEs9Aka_z9nM"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.GeocodingRequest{
		Address: "Unitec, Tegucigalpa, Honduras",
	}
	resp, err := client.Geocode(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	pretty.Print(resp)

}
