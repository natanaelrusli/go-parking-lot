package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/natanaelrusli/parking-lot/attendant"
	"github.com/natanaelrusli/parking-lot/car"
	"github.com/natanaelrusli/parking-lot/fee"
	"github.com/natanaelrusli/parking-lot/models"
	"github.com/natanaelrusli/parking-lot/parkinglot"
	parking_styles "github.com/natanaelrusli/parking-lot/parkingstyles"
)

type ParkingLotCLI struct {
	attendant attendant.ParkingAttendantItf
	tickets   map[string]*models.Ticket // Store tickets for fee calculation
}

func NewParkingLotCLI() *ParkingLotCLI {
	return &ParkingLotCLI{
		tickets: make(map[string]*models.Ticket),
	}
}

func (cli *ParkingLotCLI) Start() {
	fmt.Println("Welcome to Parking Lot System")
	fmt.Println("-----------------------------")

	// Initial setup
	cli.setupParkingLot()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nAvailable commands:")
		fmt.Println("1. Park car")
		fmt.Println("2. Unpark car")
		fmt.Println("3. Show parking lot status")
		fmt.Println("4. Change parking strategy")
		fmt.Println("5. Change fee strategy")
		fmt.Println("6. Exit")
		fmt.Print("\nEnter command number: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			cli.handleParkCar(reader)
		case "2":
			cli.handleUnparkCar(reader)
		case "3":
			cli.showStatus()
		case "4":
			cli.changeParkingStrategy(reader)
		case "5":
			cli.changeFeeStrategy(reader)
		case "6":
			fmt.Println("Thank you for using Parking Lot System")
			return
		default:
			fmt.Println("Invalid command")
		}
	}
}

func (cli *ParkingLotCLI) setupParkingLot() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter number of parking lots: ")
	input, _ := reader.ReadString('\n')
	numLots, _ := strconv.Atoi(strings.TrimSpace(input))

	parkingLots := make([]*parkinglot.ParkingLot, 0)

	for i := 0; i < numLots; i++ {
		fmt.Printf("Enter capacity for parking lot %d: ", i+1)
		input, _ := reader.ReadString('\n')
		capacity, _ := strconv.Atoi(strings.TrimSpace(input))

		lot := parkinglot.New(capacity)
		parkingLots = append(parkingLots, lot.(*parkinglot.ParkingLot))
	}

	fmt.Print("Enter attendant name: ")
	attendantName, _ := reader.ReadString('\n')
	attendantName = strings.TrimSpace(attendantName)

	cli.attendant = attendant.NewParkingAttendant(attendantName, parkingLots)
	fmt.Printf("Parking lot system initialized with %d lots and attendant %s\n", numLots, attendantName)
}

func (cli *ParkingLotCLI) handleParkCar(reader *bufio.Reader) {
	fmt.Print("Enter car license plate: ")
	licensePlate, _ := reader.ReadString('\n')
	licensePlate = strings.TrimSpace(licensePlate)

	newCar := car.NewCar(licensePlate)
	ticket, err := cli.attendant.ParkCar(newCar)

	if err != nil {
		fmt.Printf("Error parking car: %v\n", err)
		return
	}

	cli.tickets[ticket.TicketNumber] = ticket
	fmt.Printf("Car parked successfully. Ticket number: %s\n", ticket.TicketNumber)
}

func (cli *ParkingLotCLI) handleUnparkCar(reader *bufio.Reader) {
	fmt.Print("Enter ticket number: ")
	ticketNumber, _ := reader.ReadString('\n')
	ticketNumber = strings.TrimSpace(ticketNumber)

	ticket, exists := cli.tickets[ticketNumber]
	if !exists {
		fmt.Println("Invalid ticket number")
		return
	}

	unparkedCar, err := cli.attendant.UnparkCar(ticket)
	if err != nil {
		fmt.Printf("Error unparking car: %v\n", err)
		return
	}

	// Calculate parking duration and fee
	duration := time.Since(ticket.EntryTime)

	fmt.Printf("Car with license plate %s has been unparked\n", unparkedCar.LicensePlate)
	fmt.Printf("Parking duration: %v\n", duration.Round(time.Minute))

	delete(cli.tickets, ticketNumber)
}

func (cli *ParkingLotCLI) showStatus() {
	fmt.Println("\nParking Lot Status:")
	fmt.Println("-----------------")

	for _, lot := range cli.attendant.GetParkingLots() {
		fmt.Printf("Parking lot %s: %d/%d spaces occupied\n",
			lot.GetId(), lot.GetParkedCarCount(), lot.GetCapacity())
	}
}

func (cli *ParkingLotCLI) changeParkingStrategy(reader *bufio.Reader) {
	fmt.Println("\nAvailable parking strategies:")
	fmt.Println("1. Most Capacity First")
	fmt.Println("2. Most Free Space First")
	fmt.Print("Enter strategy number: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var strategy parking_styles.ParkingStyleStrategy
	switch input {
	case "1":
		strategy = parking_styles.NewMostCapacityStrategy()
		fmt.Println("Changed to Most Capacity First strategy")
	case "2":
		strategy = parking_styles.NewMostFreeSpaceStrategy()
		fmt.Println("Changed to Most Free Space First strategy")
	default:
		fmt.Println("Invalid strategy number")
		return
	}

	cli.attendant.ChangeParkingStrategy(strategy)
}

func (cli *ParkingLotCLI) changeFeeStrategy(reader *bufio.Reader) {
	fmt.Println("\nAvailable fee strategies:")
	fmt.Println("1. Hourly Rate")
	fmt.Println("2. Flat Rate")
	fmt.Print("Enter strategy number: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var strategy fee.ParkingFeeStrategy
	switch input {
	case "1":
		fmt.Print("Enter hourly rate: $")
		rateInput, _ := reader.ReadString('\n')
		rate, _ := strconv.ParseFloat(strings.TrimSpace(rateInput), 64)
		strategy = fee.NewHourlyFeeStrategy(rate)
		fmt.Printf("Changed to Hourly Rate strategy (%.2f/hour)\n", rate)
	case "2":
		fmt.Print("Enter flat rate: $")
		rateInput, _ := reader.ReadString('\n')
		rate, _ := strconv.ParseFloat(strings.TrimSpace(rateInput), 64)
		strategy = fee.NewFlatFeeStrategy(rate)
		fmt.Printf("Changed to Flat Rate strategy ($%.2f)\n", rate)
	default:
		fmt.Println("Invalid strategy number")
		return
	}

	for _, lot := range cli.attendant.GetParkingLots() {
		lot.ChangeFeeStrategy(strategy)
	}
}
