package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kr/pretty"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

const GOOGLE_GEOLOCATION_API_KEY = "AIzaSyDHjrRQfVX98LLmvEMTsjLbrQvspG4QbOc"
const GOOGLE_PLACES_API_KEY = "AIzaSyDSOkL3swVzkqHc6A99qa791zD2h4qSQpY"

//-- Routes Structs --//
type Route struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

type Coordinates struct {
	Longitude string `json:"long"`
	Latitude  string `json:"lat"`
}

//-- Places Structs --//
type Place struct {
	Origin string `json:"origin"`
}

type Image_Redux struct {
	Nombre string `json:"nombre"`
	//Data   string `json:"data"`
	Size Size `json:size`
}

type Image_GrayScale struct {
	Nombre string `json:"nombre"`
}

type Size struct {
	Alto  int `json:alto`
	Ancho int `json:ancho`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ejercicio1", CreateRoute).Methods("POST")
	router.HandleFunc("/ejercicio2", GetGeocodeDetails).Methods("POST")
	router.HandleFunc("/ejercicio3", Redux).Methods("POST")
	router.HandleFunc("/ejercicio4", GrayScaling).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetGeocodeDetails(w http.ResponseWriter, req *http.Request) {

	var route Route
	_ = json.NewDecoder(req.Body).Decode(&route)

	client, err := maps.NewClient(maps.WithAPIKey(GOOGLE_GEOLOCATION_API_KEY))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := &maps.GeocodingRequest{
		Address: route.Origin,
	}

	r2 := &maps.NearbySearchRequest{
		Radius: 500,
		Type:   maps.PlaceTypeRestaurant,
	}

	ll, err := maps.ParseLatLng(route.Origin)
	pretty.Print(ll.Lat)

	resp, err := client.Geocode(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	resp2, err := client.NearbySearch(context.Background(), r2)
	pretty.Print(resp)
	pretty.Println(resp2)
	//pretty.Print(r)
	//json.NewEncoder(w).Encode(resp)
	//pretty.Println(x)
}

func CreateRoute(w http.ResponseWriter, req *http.Request) {
	//Request Origin and Destination.
	var route Route
	_ = json.NewDecoder(req.Body).Decode(&route)

	//Use google directions api with requested origin and destination.
	client, err := maps.NewClient(maps.WithAPIKey(GOOGLE_GEOLOCATION_API_KEY))
	if err != nil {
		log.Fatalf("Fatal Error: %s", err)
	}

	r := &maps.DirectionsRequest{

		Origin:      route.Origin,
		Destination: route.Destination,
		Mode:        maps.TravelModeDriving,
	}

	routes, _, err := client.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("Fatal Error: %s", err)
	}

	//pretty.Println(routes)
	//pretty.Println(location)
	json.NewEncoder(w).Encode(routes)
}

/*
func CreatePlacesList(w http.ResponseWriter, req *http.Request) {
    //Request Origin and Destination.
    var place Place
    _ = json.NewDecoder(req.Body).Decode(&place)
    json.NewEncoder(w).Encode(place)

    client, err := maps.NewClient(maps.WithAPIKey(GOOGLE_PLACES_API_KEY))
    if err != nil {
        log.Fatalf("Fatal Error: %s", err)
    }
    r := &maps.NearbySearchRequest{

        Radius:  80,
        Keyword: "restaurant",
        Name:    place.Origin,
    }

    places, _, err := client.NearbySearch(context.Background(), r)
    if err != nil {
        log.Fatalf("Fatal Error: %s", err)
    }

    pretty.Println(places)
}
*/

func Redux(w http.ResponseWriter, req *http.Request) {

	var img Image_Redux
	_ = json.NewDecoder(req.Body).Decode(&img)

	bitmap, err := openImage(img.Nombre)
	if err != nil {
		fmt.Println(err)
	}

	pretty.Println("YASS")

	bounds := bitmap.Bounds()
	width, height := bounds.Max.X/img.Size.Alto, bounds.Max.Y/img.Size.Ancho

	imgSet := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < img.Size.Alto; y++ {
		for x := 0; x < img.Size.Ancho; x++ {
			pixel := bitmap.At(x*width, y*height)
			imgSet.Set(x, y, pixel)
		}
	}

	outfile, err := os.Create("lena_Redux.bmp")
	if err != nil {
		fmt.Println(err)
	}

	defer outfile.Close()

	//json.NewEncoder(w).Encode(imgSet)
	pretty.Println(imgSet)
	bmp.Encode(outfile, imgSet)

	res, err := openImage("lena_Redux.bmp")
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(res.ColorModel())
	pretty.Println(res)
	//json.NewEncoder(w).Encode(bmp)

}

func GrayScaling(w http.ResponseWriter, req *http.Request) {

	var img Image_GrayScale
	_ = json.NewDecoder(req.Body).Decode(&img)

	bitmap, err := openImage(img.Nombre)
	if err != nil {
		fmt.Println(err)
	}

	bounds := bitmap.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	imgSet := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			oldPixel := bitmap.At(x, y)
			r, g, b, _ := oldPixel.RGBA()

			avg := (r + g + b) / 3
			pixel := color.Gray{uint8(avg / 256)}

			imgSet.Set(x, y, pixel)
		}
	}

	outfile, err := os.Create("lena_GrayScale.bmp")
	if err != nil {
		fmt.Println(err)
	}

	defer outfile.Close()

	//json.NewEncoder(w).Encode(imgSet)

	pretty.Println(imgSet)

	bmp.Encode(outfile, imgSet)
}

func openImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return bmp.Decode(f)
}
