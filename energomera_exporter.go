package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/peak-load/energomera"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tarm/serial"
	"github.com/tkanos/gonfig"
)
// listen-address sets IP address and port for exporter
var addr = flag.String("listen-address", ":9876", "The address to listen on for HTTP requests.")

// config sets path to JSON file with energomera settings
var config = flag.String("config", "config.json", "Config file path for energomera.")

// Configuration for connection
type Configuration struct {
	Port          string
	SleepInterval time.Duration
	Counters      []string
}

func main() {
	flag.Parse()

	r := prometheus.NewRegistry()

	voltage := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltage_volt",
			Help: "voltage per phase in volts",
		},
		[]string{"id", "phase"},
	)

	current := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "current_ampere",
			Help: "current per phase in amperes",
		},
		[]string{"id", "phase"},
	)

	powerp := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "power_kwt",
			Help: "power per phase in kilowatts",
		},
		[]string{"id", "phase"},
	)

	tarif := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tarif",
			Help: "tarif total in kilowatts",
		},
		[]string{"id", "tarif"},
	)

	power := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "power_used_kwt",
			Help: "power used total in kilowatts",
		},
		[]string{"id"},
	)

	frequency := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mains_frequency_hz",
			Help: "mains frequency in Hertz",
		},
		[]string{"id"},
	)

	r.MustRegister(voltage)
	r.MustRegister(current)
	r.MustRegister(powerp)
	r.MustRegister(tarif)
	r.MustRegister(power)
	r.MustRegister(frequency)

	configuration := Configuration{}

	err := gonfig.GetConf(*config, &configuration)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 128)

	port := "/dev/ttyS0"
	counters := []string{""}
	SleepInterval := time.Millisecond * 500

	if len(configuration.Port) > 0 {
		port = configuration.Port
	}

	if len(configuration.Counters) > 0 {
		counters = configuration.Counters
	}

	SleepInterval = configuration.SleepInterval
	log.Printf("Data collection started")

	go func() {
		for {
			start := time.Now()
			for counter := range counters {

				c := &serial.Config{Name: port, Baud: 9600, Size: 7, StopBits: 1, Parity: 'E', ReadTimeout: time.Millisecond * SleepInterval}

				time.Sleep(time.Second * 1)

				s, err := serial.OpenPort(c)

				if err != nil {
					log.Fatal(err)
				}
				// initialize commands
				n, _ := s.Write([]byte("/?" + counters[counter] + "!\r\n"))
				time.Sleep(time.Millisecond * SleepInterval)
				n, _ = s.Read(buf)

				// clounter identifier
				ident := strings.Split(string(buf[:n]), "")
				time.Sleep(time.Millisecond * SleepInterval)

				n, _ = s.Write([]byte("\x060" + ident[4] + "1\r\n"))
				time.Sleep(time.Millisecond * SleepInterval)
				n, _ = s.Read(buf)

				// send commands
				commands := []string{"VOLTA", "CURRE", "POWEP", "POWPP", "FREQU", "ET0PE"}
				commandline := map[string]string{"head": "R1", "body": ""}

				for i := range commands {
					commandline["body"] = commands[i] + "()"
					command := energomera.DataEncode(commandline)
					n, _ = s.Write(command)
					time.Sleep(time.Millisecond * SleepInterval)
					n, _ = s.Read(buf)

					switch commands[i] {

					case "VOLTA":
						phasev := strings.Split(string(buf[1:n]), "\r\n")

						phase1v, _ := strconv.ParseFloat(strings.Trim(phasev[0], "VOLTA()"), 8)
						voltage.WithLabelValues(counters[counter], "phase1").Set(phase1v)

						phase2v, _ := strconv.ParseFloat(strings.Trim(phasev[1], "VOLTA()"), 8)
						voltage.WithLabelValues(counters[counter], "phase2").Set(phase2v)

						phase3v, _ := strconv.ParseFloat(strings.Trim(phasev[2], "VOLTA()"), 8)
						voltage.WithLabelValues(counters[counter], "phase3").Set(phase3v)

					case "CURRE":
						phasea := strings.Split(string(buf[1:n]), "\r\n")

						phase1a, _ := strconv.ParseFloat(strings.Trim(phasea[0], "CURRE()"), 8)
						current.WithLabelValues(counters[counter], "phase1").Set(phase1a)

						phase2a, _ := strconv.ParseFloat(strings.Trim(phasea[1], "CURRE()"), 8)
						current.WithLabelValues(counters[counter], "phase2").Set(phase2a)

						phase3a, _ := strconv.ParseFloat(strings.Trim(phasea[2], "CURRE()"), 8)
						current.WithLabelValues(counters[counter], "phase3").Set(phase3a)

					case "POWPP":
						phasep := strings.Split(string(buf[1:n]), "\r\n")

						phase1p, _ := strconv.ParseFloat(strings.Trim(phasep[0], "POWPP()"), 8)
						powerp.WithLabelValues(counters[counter], "phase1").Set(phase1p)

						phase2p, _ := strconv.ParseFloat(strings.Trim(phasep[1], "POWPP()"), 8)
						powerp.WithLabelValues(counters[counter], "phase2").Set(phase2p)

						phase3p, _ := strconv.ParseFloat(strings.Trim(phasep[2], "POWPP()"), 8)
						powerp.WithLabelValues(counters[counter], "phase3").Set(phase3p)

					case "ET0PE":
						tarift := strings.Split(string(buf[1:n]), "\r\n")

						tarif1t, _ := strconv.ParseFloat(strings.Trim(tarift[0], "ET0PE()"), 8)
						tarif.WithLabelValues(counters[counter], "tarif1").Set(tarif1t)

						tarif2t, _ := strconv.ParseFloat(strings.Trim(tarift[1], "ET0PE()"), 8)
						tarif.WithLabelValues(counters[counter], "tarif2").Set(tarif2t)

						tarif3t, _ := strconv.ParseFloat(strings.Trim(tarift[2], "ET0PE()"), 8)
						tarif.WithLabelValues(counters[counter], "tarif3").Set(tarif3t)

					case "POWEP":
						powert := strings.Split(string(buf[1:n]), "\r\n")

						pow, _ := strconv.ParseFloat(strings.Trim(powert[0], "POWEP()"), 8)
						power.WithLabelValues(counters[counter]).Set(pow)

					case "FREQU":
						freq := strings.Split(string(buf[1:n]), "\r\n")

						freqt, _ := strconv.ParseFloat(strings.Trim(freq[0], "FREQU()"), 8)
						frequency.WithLabelValues(counters[counter]).Set(freqt)
					}
				}
				end := []byte("\x01\x42\x30\x03\x75")
				n, _ = s.Write(end)
				s.Close()
			}
			elapsed := time.Since(start)
			log.Printf("Data fetch took %s", elapsed)
			time.Sleep(time.Second * 5)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(
		r,
		promhttp.HandlerOpts{},
	))

	log.Printf("Starting web server at %s\n", *addr)
	httperr := http.ListenAndServe(*addr, nil)
	if httperr != nil {
		log.Printf("http.ListenAndServer: %v\n", httperr)
	}
}
