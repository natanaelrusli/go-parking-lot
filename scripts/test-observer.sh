#!/bin/bash

# Navigate to the parkinglot directory
cd parkinglot

# Run only the observer tests
go test -v -run "TestParkingLotObserver" ./...

# Return to original directory
cd ..