package main

import (
	"log"
	"os"
)

func Label() (input []string, output []string) {
	input = []string{"filea", "fileb"}
	output = []string{"result"}
	return
}

func Main(inputs []string, outputs []string) int {
	a, _ := os.ReadFile(inputs[0])
	b, _ := os.ReadFile(inputs[1])
	f, err := os.OpenFile(outputs[0], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return 1
	}
	defer f.Close()

	log.SetOutput(f)
	if string(a) == string(b) {
		return 0
	}
	log.Println("a != b")
	return 1
}
