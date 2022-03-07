package main

import (
	"hypercurious/test/db"
)

func main() {

	defer db.Close()
}
