package main

// The server for the mtr-ui, starting as simple as possible.
// TODO: keep 12 factor app principles in mind, use logentries, env vars instead of json config, etc.

import (
	_ "github.com/GeoNet/log/logentries"
	"github.com/GeoNet/weft"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	mux           *http.ServeMux
	mtrApiUrl     *url.URL
	webServerPort string
)

func init() {
	var err error
	webServerPort = os.Getenv("MTR_UI_PORT")
	mtrApiUrlString := os.Getenv("MTR_API_URL")
	switch "" {
	case webServerPort:
		log.Fatal("error, environment variable MTR_UI_PORT must be set (eg: 8080)")
	case mtrApiUrlString:
		log.Fatal("error, environment variable MTR_API_URL must be set (eg: https://mtr-api.geonet.org.nz)")
	}

	if mtrApiUrl, err = url.Parse(mtrApiUrlString); err != nil {
		log.Fatal(err)
	}

	mux = http.NewServeMux()
	mux.HandleFunc("/", weft.MakeHandlerPage(homePageHandler))
	mux.HandleFunc("/field", weft.MakeHandlerPage(fieldPageHandler))
	mux.HandleFunc("/field/", weft.MakeHandlerPage(fieldPageHandler))
	mux.HandleFunc("/field/metric", weft.MakeHandlerPage(metricDetailHandler))
	mux.HandleFunc("/field/metrics", weft.MakeHandlerPage(fieldMetricsPageHandler))
	mux.HandleFunc("/field/devices", weft.MakeHandlerPage(fieldDevicesPageHandler))
	mux.HandleFunc("/field/plot", weft.MakeHandlerPage(fieldPlotPageHandler))
	mux.HandleFunc("/data", weft.MakeHandlerPage(dataPageHandler))
	mux.HandleFunc("/data/", weft.MakeHandlerPage(dataPageHandler))
	mux.HandleFunc("/data/sites", weft.MakeHandlerPage(dataSitesPageHandler))
	mux.HandleFunc("/data/metrics", weft.MakeHandlerPage(dataMetricsPageHandler))
	mux.HandleFunc("/data/plot", weft.MakeHandlerPage(dataPlotPageHandler))
	mux.HandleFunc("/map", weft.MakeHandlerPage(mapPageHandler))
	mux.HandleFunc("/map/", weft.MakeHandlerPage(mapPageHandler))
	mux.HandleFunc("/search", weft.MakeHandlerPage(searchPageHandler))
	mux.HandleFunc("/tag", weft.MakeHandlerPage(tagPageHandler))
	mux.HandleFunc("/tag/", weft.MakeHandlerPage(tagPageHandler))
	mux.HandleFunc("/app", weft.MakeHandlerPage(appPageHandler))
	mux.HandleFunc("/app/", weft.MakeHandlerPage(appPageHandler))
	mux.HandleFunc("/app/plot", weft.MakeHandlerPage(appPlotPageHandler))
}

func main() {
	log.Println("starting mtr-ui server")
	log.Fatal(http.ListenAndServe(":"+webServerPort, mux))
}
