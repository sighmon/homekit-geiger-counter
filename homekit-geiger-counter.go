package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"

	"github.com/sighmon/homekit-geiger-counter/geigercounter"

	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

var sensorHost string
var sensorPort int
var secondsBetweenReadings time.Duration
var developmentMode bool
const tubeMultiplier = 6.6086956522

func init() {
	flag.StringVar(&sensorHost, "host", "http://0.0.0.0", "sensor host, a string")
	flag.IntVar(&sensorPort, "port", 1006, "sensor port number, an int")
	flag.DurationVar(&secondsBetweenReadings, "sleep", 3*time.Second, "how many seconds between sensor readings, an int followed by the duration")
	flag.BoolVar(&developmentMode, "dev", false, "turn on development mode to return a random temperature reading, boolean")
	flag.Parse()

	if developmentMode == true {
		log.Println("Development mode on, ignoring sensor and returning random values...")
	}
}

func main() {
	info := accessory.Info{
		Name:             "Radiation",
		SerialNumber:     "SEN0463",
		Manufacturer:     "DF Robot",
		Model:            "Gravity Geiger Counter",
		FirmwareRevision: "1.0.0",
	}

	acc := geigercounter.NewAccessory(
		info,
	)

	config := hc.Config{
		// Change the default Apple Accessory Pin if you wish
		Pin: "00102003",
		// Port: "12345",
		// StoragePath: "./db",
	}

	t, err := hc.NewIPTransport(config, acc.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	// Get the sensor readings every secondsBetweenReadings
	go func() {
		type Reading struct {
			Name  string
			Value float64
		}

		type Readings struct {
			Cpm     Reading
			Nsvh    Reading
			Usvh    Reading
		}

		readings := Readings{
			Cpm: Reading{
				Name:  "cpm",
				Value: 0,
			},
			Nsvh: Reading{
				Name:  "radiation_nsvh",
				Value: 0,
			},
			Usvh: Reading{
				Name:  "radiation_usvh",
				Value: 0,
			},
		}
		values := reflect.ValueOf(readings)

		for {
			// Get readings from the Prometheus exporter
			resp, err := http.Get(fmt.Sprintf("%s:%d", sensorHost, sensorPort))
			if err == nil {
				defer resp.Body.Close()
				scanner := bufio.NewScanner(resp.Body)
				for scanner.Scan() {
					line := scanner.Text()
					// Parse the readings
					for i := 0; i < values.NumField(); i++ {
						fieldname := values.Field(i).Interface().(Reading).Name
						regexString := fmt.Sprintf("^%s", fieldname) + ` ([-+]?\d*\.\d+|\d+)`
						re := regexp.MustCompile(regexString)
						rs := re.FindStringSubmatch(line)
						if rs != nil {
							parsedValue, err := strconv.ParseFloat(rs[1], 64)
							if err == nil {
								if developmentMode {
									println(fmt.Sprintf("%s %f", fieldname, parsedValue))
								}
								switch fieldname {
								case "cpm":
									readings.Cpm.Value = parsedValue
								case "radiation_nsvh":
									readings.Nsvh.Value = parsedValue
								case "radiation_usvh":
									readings.Usvh.Value = parsedValue
								}
							}
						}
					}
				}
				scanner = nil
			} else {
				log.Println(err)
			}

			if developmentMode {
				// Return a random float between 25 and 100
				randomValue := 25 + rand.Float64()*(100-25)
				readings.Cpm.Value = randomValue
				readings.Nsvh.Value = randomValue*tubeMultiplier
				readings.Usvh.Value = (randomValue*tubeMultiplier)/1000
			}

			// Set the sensor readings
			acc.GeigerCounter.Cpm.SetValue(readings.Cpm.Value)
			acc.GeigerCounter.Nsvh.SetValue(readings.Nsvh.Value)
			acc.GeigerCounter.Usvh.SetValue(readings.Usvh.Value)

			log.Println(fmt.Sprintf("CPM: %f", readings.Cpm.Value))
			log.Println(fmt.Sprintf("nSv/h: %f", readings.Nsvh.Value))
			log.Println(fmt.Sprintf("uSv/h: %f", readings.Usvh.Value))

			// Time between readings
			time.Sleep(secondsBetweenReadings)
		}
	}()

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
