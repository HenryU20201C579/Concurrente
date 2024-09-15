// Segundo intento
package main

import (
	"fmt"
	"time"
)

var wantp bool = false
var wantq bool = false

func p() {
	for {
		fmt.Println("Line1-SNC P")
		for wantq != false {
			//espera
		}
		wantp = true
		fmt.Println("Line1-SC P")
		wantp = false
	}
}

func q() {
	for {
		fmt.Println("Line1-SNC Q")
		for wantp != false {
			//espera
		}
		wantq = true
		fmt.Println("Line1-SC Q")
		wantq = false
	}
}

func main() {
	go p()
	go q()
	
	time.Sleep(time.Hour)
}
