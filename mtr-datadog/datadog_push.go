package main

// Test program that pulls metric data from mtr-api and pushed to DataDog using their statsd package

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/GeoNet/mtr/mtrpb"
	"github.com/golang/protobuf/proto"
	"github.com/lunny/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const SAMPLE_URL = "https://mtr-api.geonet.org.nz/data/latency?siteID=CECS&typeID=latency.strong"

func getBytes(urlString string, accept string) (body []byte, err error) {
	var client = &http.Client{}
	var request *http.Request
	var response *http.Response

	if request, err = http.NewRequest("GET", urlString, nil); err != nil {
		return nil, err
	}
	//request.SetBasicAuth(userW, keyW)
	request.Header.Add("Accept", accept)

	if response, err = client.Do(request); err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		msg := ""
		if response.Body != nil {
			if b, err := ioutil.ReadAll(response.Body); err == nil {
				msg = ":" + string(b)
			}
		}

		return nil, fmt.Errorf("Wrong response code for %s got %d expected %d %s", urlString, response.StatusCode, http.StatusOK, msg)
	}

	// Read body, could use io.LimitReader() to avoid a massive read (unlikely)
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getMetrics() (*mtrpb.DataLatencyResult, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(SAMPLE_URL); err != nil {
		return nil, err
	}

	var protoData []byte
	if protoData, err = getBytes(u.String(), "application/x-protobuf"); err != nil {
		return nil, err
	}

	var f mtrpb.DataLatencyResult
	if err = proto.Unmarshal(protoData, &f); err != nil {
		return nil, err
	}

	return &f, nil
}

// push metrics to the DataDog API and see what happens. Example from https://godoc.org/github.com/DataDog/datadog-go/statsd:
// Can't control the timestamp of metrics using statsd package, but could do it using the API: http://docs.datadoghq.com/api/
func dataDogPush(data *mtrpb.DataLatencyResult) error {
	// Create the client
	//c, err := statsd.NewBuffered("127.0.0.1:8125", 100)
	c, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Prefix every metric with the app name
	c.Namespace = "test."
	c.Tags = append(c.Tags, "strong-motion")

	i := 0
	for {
		i = i % len(data.Result)
		fmt.Println("About to send a metric", data.Result[i].Mean)
		if err = c.Gauge("latency", float64(data.Result[i].Mean), nil, 1); err != nil {
			return err
		}
		fmt.Println("  sent a metric!")
		time.Sleep(time.Minute)
		i++
	}

	return nil
}

func main() {
	var data *mtrpb.DataLatencyResult
	var err error
	if data, err = getMetrics(); err != nil {
		log.Fatal(err)
	}

	if err = dataDogPush(data); err != nil {
		log.Fatal(err)
	}
}
