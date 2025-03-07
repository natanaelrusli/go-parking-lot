package main

import (
	"github.com/natanaelrusli/parking-lot/cli"
)

func main() {
	parkingLotCLI := cli.NewParkingLotCLI()
	parkingLotCLI.Start()
}
