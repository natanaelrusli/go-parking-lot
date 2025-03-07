# Run observer tests
test:
	go test -v ./...

make-mock:
	mockery --name ParkingLotObserver --dir models --output mocks
