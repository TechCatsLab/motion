package main

import (
	"fmt"

	"github.com/TechCatsLab/motion/ssh-zh/conn"
)

func main() {
	ct:= conn.Client
	defer ct.Close()

	session, err := ct.NewSession()
	var outBytes []byte
	outBytes, err = session.Output("echo hello")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(outBytes))
	session.Close()
}
