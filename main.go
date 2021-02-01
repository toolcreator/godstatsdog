package main

import (
	"log"
)

func main() {
	stats, err := getDStats()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(stats)
	}
}
