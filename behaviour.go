package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type behaviour struct {
	Delay     time.Duration
	Frequency float64
	Status    int
	Body      string
	Garbage   int
}

type randomiser interface {
	getFloat() float64
}

type defaultRandomiser struct{}

func (d *defaultRandomiser) getFloat() float64 {
	return rand.Float64()
}

func getBehaviour(behaviours []behaviour, randomiser randomiser) *behaviour {
	randnum := randomiser.getFloat()
	lower := 0.0
	var upper float64
	for _, behaviour := range behaviours {
		upper = lower + behaviour.Frequency
		if randnum > lower && randnum <= upper {
			return &behaviour
		}

		lower = upper
	}

	return nil

}

func (b behaviour) String() string {

	frequency := fmt.Sprintf("%2.0f%% of the time |", b.Frequency*100)

	delay := ""
	if b.Delay != 0 {
		delay = fmt.Sprintf("Delay: %v ", b.Delay*time.Millisecond)
	}

	status := ""
	if b.Status != 0 {
		status = fmt.Sprintf("Status: %v ", b.Status)
	}

	body := ""
	if b.Body != "" {
		body = fmt.Sprintf("Body: %v ", b.Body)
	}

	garbage := ""
	if b.Garbage != 0 {
		garbage = fmt.Sprintf("Garbage bytes: %d ", b.Garbage)
	}

	return fmt.Sprintf("%v %v%v%v%v", frequency, delay, status, body, garbage)
}

func loadMonkeyConfig(path string) []behaviour {

	if path == "" {
		return []behaviour{}
	}

	config, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalf("Problem occured when trying to read the config file: %v", err)
	}

	behaviours := monkeyConfigFromYAML(config)

	log.Println("Monkey config loaded")
	for _, b := range behaviours {
		log.Println(b)
	}

	return behaviours
}

func monkeyConfigFromYAML(data []byte) []behaviour {
	var result []behaviour
	err := yaml.Unmarshal([]byte(data), &result)

	if err != nil {
		log.Fatalf("Problem occured when trying to parse the config file: %v", err)
	}

	return result
}