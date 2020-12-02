package main

import (
	"farmtotable/aragorn"
	"farmtotable/legolas"
)

/* Backend backend_launcher. This can be used for dev/test purposes when we want both
Aragorn and Legolas running as the same service allowing us to use Sqlite for
Gandalf's backend. THIS MUST NOT BE USED IN PRODUCTION. */
func main() {
	go aragorn.NewAragorn().Run()
	go legolas.NewLegolas().Run()
	// Block forever
	select {}
}
