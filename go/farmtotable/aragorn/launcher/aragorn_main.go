package main

import (
	"farmtotable/aragorn"
	"flag"
)

/* Aragorn backend_launcher */
func main() {
	flag.Parse()
	aragorn.NewAragorn().Run()
}
