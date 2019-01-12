package main

import (
	"log"

	"github.com/TechCatsLab/motion/deploy"
)

func main() {
	log.Fatalln(deploy.Start())
}
