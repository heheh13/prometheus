package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)
var (
	appVersion string
	version = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "version",
		Help:        "version infomation about this binary",
		ConstLabels: map[string]string{
			"version":appVersion,
		},
	})
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "http_requests_total",
		Help:        "count of all http request",
	},[]string{"code","method"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "http_request_duration_second",
		Help:        "Duration of all Http request",
	},[]string{"code","handler","method"})
)

func notfound(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found form api server"))
}
func home (writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("hello from api server"))
}
func main() {
	bind := ""
	flagset := flag.NewFlagSet(os.Args[0],flag.ExitOnError)
	flagset.StringVar(&bind,"bind",":8080","the socket to bind to")
	flagset.Parse(os.Args[1:])

	r := prometheus.NewRegistry()
	r.MustRegister(httpRequestsTotal)
	r.MustRegister(httpRequestDuration)
	r.MustRegister(version)

	handler := mux.NewRouter()

	homeHandler := http.HandlerFunc(home)

	notFoundHandler := http.HandlerFunc(notfound)

	handler.HandleFunc("/",promhttp.InstrumentHandlerCounter(httpRequestsTotal,homeHandler))
	handler.HandleFunc("/err",promhttp.InstrumentHandlerCounter(httpRequestsTotal,notFoundHandler))

	handler.Handle("/metrics",promhttp.HandlerFor(r,promhttp.HandlerOpts{}))

	fmt.Printf("server is running on port %s\n",bind)
	log.Fatal(http.ListenAndServe(bind,handler))
}
