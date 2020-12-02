package main

import (
	"farmtotable/legolas"
	"flag"
)

/* Legolas backend_launcher */
func main() {
	flag.Parse()
	legolas.NewLegolas().Run()
}
