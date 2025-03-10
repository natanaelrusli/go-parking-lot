# Run observer tests
test:
	go test -v ./...

make-mock:
    mockery --name ParkingLotObserver --dir models --output mocks
    mockery --name ParkingFeeStrategy --dir fee --output mocks
