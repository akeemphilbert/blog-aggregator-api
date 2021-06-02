package main

import (
	"flag"

	api "github.com/wepala/blog-aggregator-api/src"
)

var port = flag.String("port","8682","-port=8682")

func main() {
	api.New(port,"")
}