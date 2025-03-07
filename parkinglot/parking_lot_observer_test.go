package parkinglot

import (
	"log"
	"testing"

	"github.com/natanaelrusli/parking-lot/car"
	"github.com/natanaelrusli/parking-lot/mocks"
	"github.com/natanaelrusli/parking-lot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockObserver implements ParkingLotObserver for testing
type MockObserver struct {
	notifications []models.ParkingLotStatus
	name          string
}

type LoggingObserver struct {
	logFile string
}

func NewLoggingObserver(logFile string) *LoggingObserver {
	return &LoggingObserver{
		logFile: logFile,
	}
}

func NewMockObserver(name string) *MockObserver {
	return &MockObserver{
		notifications: make([]models.ParkingLotStatus, 0),
		name:          name,
	}
}

func (lo *LoggingObserver) OnParkingLotStatusChanged(status models.ParkingLotStatus) {
	log.Printf("Logging to file %s: %v", lo.logFile, status)
}

func (m *MockObserver) OnParkingLotStatusChanged(status models.ParkingLotStatus) {
	m.notifications = append(m.notifications, status)
}

func TestParkingLotObserverNotifications(t *testing.T) {
	t.Run("should notify observer when a car is parked", func(t *testing.T) {
		// Arrange
		parkingLot := New(2)
		observer := NewMockObserver("TestObserver")
		// adding observer to parking lot
		parkingLot.AddObserver(observer)

		// Act
		car := car.NewCar("ABC123")
		_, err := parkingLot.Park(car)

		// Assert
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(observer.notifications) != 1 {
			t.Errorf("Expected 1 notification, got %d", len(observer.notifications))
		}

		lastNotification := observer.notifications[0]
		if lastNotification.Available != 1 {
			t.Errorf("Expected 1 available space, got %d", lastNotification.Available)
		}
		if lastNotification.IsFull {
			t.Error("Parking lot should not be full")
		}
	})

	t.Run("should notify observer when parking lot becomes full", func(t *testing.T) {
		// Arrange
		parkingLot := New(2)
		observer := NewMockObserver("TestObserver")
		parkingLot.AddObserver(observer)

		// Act
		car1 := car.NewCar("ABC123")
		car2 := car.NewCar("XYZ789")
		_, _ = parkingLot.Park(car1)
		_, _ = parkingLot.Park(car2)

		// Assert
		if len(observer.notifications) != 2 {
			t.Errorf("Expected 2 notifications, got %d", len(observer.notifications))
		}

		lastNotification := observer.notifications[len(observer.notifications)-1]
		if !lastNotification.IsFull {
			t.Error("Parking lot should be full")
		}
		if lastNotification.Available != 0 {
			t.Errorf("Expected 0 available spaces, got %d", lastNotification.Available)
		}
	})

	t.Run("should not notify a non observer", func(t *testing.T) {
		// Arrange
		parkingLot := New(2)
		observer := NewMockObserver("TestObserver")

		// Act
		car := car.NewCar("ABC123")
		parkingLot.Park(car)

		// Assert
		if len(observer.notifications) != 0 {
			t.Errorf("Expected 0 notifications, got %d", len(observer.notifications))
		}
	})

	t.Run("should notify observer when car is unparked", func(t *testing.T) {
		// Arrange
		parkingLot := New(2)
		observer := NewMockObserver("TestObserver")
		parkingLot.AddObserver(observer)

		// Park and then unpark a car
		car := car.NewCar("ABC123")
		ticket, _ := parkingLot.Park(car)
		_, _ = parkingLot.Unpark(ticket)

		// Assert
		if len(observer.notifications) != 2 {
			t.Errorf("Expected 2 notifications, got %d", len(observer.notifications))
		}

		lastNotification := observer.notifications[len(observer.notifications)-1]
		if lastNotification.Available != 2 {
			t.Errorf("Expected 2 available spaces, got %d", lastNotification.Available)
		}
		if lastNotification.IsFull {
			t.Error("Parking lot should not be full")
		}
	})

	t.Run("should support multiple observers", func(t *testing.T) {
		// Arrange
		parkingLot := New(2)
		observer1 := NewMockObserver("Observer1")
		observer2 := NewMockObserver("Observer2")
		parkingLot.AddObserver(observer1)
		parkingLot.AddObserver(observer2)

		// Act
		car := car.NewCar("ABC123")
		_, _ = parkingLot.Park(car)

		// Assert
		if len(observer1.notifications) != 1 {
			t.Errorf("Observer1: Expected 1 notification, got %d", len(observer1.notifications))
		}
		if len(observer2.notifications) != 1 {
			t.Errorf("Observer2: Expected 1 notification, got %d", len(observer2.notifications))
		}

		// Verify both observers received the same notification
		notification1 := observer1.notifications[0]
		notification2 := observer2.notifications[0]
		if notification1.Available != notification2.Available {
			t.Error("Observers received different notifications")
		}
	})

	t.Run("should include correct lot ID in notifications", func(t *testing.T) {
		// Arrange
		parkingLot := New(2)
		observer := NewMockObserver("TestObserver")
		parkingLot.AddObserver(observer)

		// Act
		car := car.NewCar("ABC123")
		_, _ = parkingLot.Park(car)

		// Assert
		if len(observer.notifications) != 1 {
			t.Errorf("Expected 1 notification, got %d", len(observer.notifications))
		}

		notification := observer.notifications[0]
		if notification.LotID != parkingLot.ID {
			t.Errorf("Expected lot ID %s, got %s", parkingLot.ID, notification.LotID)
		}
	})

	t.Run("should notify with correct capacity values", func(t *testing.T) {
		// Arrange
		capacity := 3
		parkingLot := New(capacity)
		observer := NewMockObserver("TestObserver")
		parkingLot.AddObserver(observer)

		// Act
		car := car.NewCar("ABC123")
		parkingLot.Park(car)

		// Assert
		if len(observer.notifications) != 1 {
			t.Errorf("Expected 1 notification, got %d", len(observer.notifications))
		}

		notification := observer.notifications[0]
		if notification.Capacity != capacity {
			t.Errorf("Expected capacity %d, got %d", capacity, notification.Capacity)
		}

		if notification.Available != capacity-1 {
			t.Errorf("Expected available spaces %d, got %d", capacity-1, notification.Available)
		}
	})

	t.Run("should notify only subscribers", func(t *testing.T) {
		// arrange
		parkingLot := New(2)
		observer1 := NewMockObserver("Observer1")
		observer2 := NewMockObserver("Observer2")
		parkingLot.AddObserver(observer1)

		// act
		car := car.NewCar("ABC123")
		parkingLot.Park(car)

		log.Println("obs1", observer1.notifications)
		log.Println("obs2", observer2.notifications)

		// assert
		if len(observer1.notifications) != 1 {
			t.Errorf("Observer1: Expected 1 notification, got %d", len(observer1.notifications))
		}

		if len(observer2.notifications) != 0 {
			t.Errorf("Observer2: Expected 0 notifications, got %d", len(observer2.notifications))
		}
	})
}

func TestParkingLotLoggingObserver(t *testing.T) {
	t.Run("should support multiple types of observers", func(t *testing.T) {
		parkingLot := New(2)
		mockObs := NewMockObserver("MockObserver")
		logObs := NewLoggingObserver("parking.log")

		parkingLot.AddObserver(mockObs)
		parkingLot.AddObserver(logObs)

		car := car.NewCar("ABC123")
		_, _ = parkingLot.Park(car)

		if len(mockObs.notifications) != 1 {
			t.Errorf("Expected 1 notification for mock observer, got %d", len(mockObs.notifications))
		}

		if len(logObs.logFile) == 0 {
			t.Error("Expected log file to be set for logging observer")
		}
	})
}

func TestParkingLotObserverWithMockery(t *testing.T) {
	t.Run("should notify observer when car is parked", func(t *testing.T) {
		// arrange
		parkingLot := New(2)
		mockObserver := mocks.NewParkingLotObserver(t)

		mockObserver.On("OnParkingLotStatusChanged", mock.MatchedBy(func(status models.ParkingLotStatus) bool {
			return status.Available == 1 && !status.IsFull
		})).Return(nil).Once()

		// add observer and perform action
		parkingLot.AddObserver(mockObserver)
		car := car.NewCar("ABC123")
		ticket, err := parkingLot.Park(car)

		// assert
		assert.NotNil(t, ticket)
		assert.NoError(t, err)
		mockObserver.AssertExpectations(t)
	})

	t.Run("should notify observer when parking lot becomes full", func(t *testing.T) {
		// arrange
		pl1 := New(1)
		mockObserver := mocks.NewParkingLotObserver(t)
		mockObserver.On("OnParkingLotStatusChanged", mock.MatchedBy(func(status models.ParkingLotStatus) bool {
			return status.IsFull
		})).Return(nil).Once()

		// add observer and perform action to trigger notify
		pl1.AddObserver(mockObserver)
		c1 := car.NewCar("B2534POZ")
		ticket, err := pl1.Park(c1)

		// assert
		assert.NotNil(t, ticket)
		assert.NoError(t, err)

		// mockObserver.AssertExpectations(t)
	})
}
