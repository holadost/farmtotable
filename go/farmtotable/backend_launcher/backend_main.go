package main

import (
	"farmtotable/aragorn"
	"farmtotable/gandalf"
	"farmtotable/legolas"
	"flag"
)

/* Backend backend_launcher. This can be used for dev/test purposes when we want both
Aragorn and Legolas running as the same service allowing us to use Sqlite for
Gandalf's backend. THIS MUST NOT BE USED IN PRODUCTION. */
func main() {
	flag.Parse()
	g := gandalf.NewSqliteGandalf()
	go aragorn.NewAragornWithGandalf(g).Run()
	go legolas.NewLegolasWithGandalf(g).Run()
	// Block forever
	select {}
}
