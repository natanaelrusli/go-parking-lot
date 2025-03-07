# go-parking-lot

## Observer Pattern

### Flow of Information:

#### When a car is parked or unparked:
- The parking lot's state changes
notifyObservers() is called
- A ParkingLotStatus object is created with current state
- Each observer receives this status via OnParkingLotStatusChanged()

#### Observers can:
- Track parking lot capacity
- Monitor when lots become full
- Log parking events
- Trigger actions based on lot status
- Display status updates

#### This implementation allows for:
- Multiple observers watching the same parking lot
- Different types of observers (logging, monitoring, UI updates, etc.)
- Easy addition/removal of observers
- Consistent status updates across all observers

#### The pattern is particularly useful here because:
- It decouples the parking lot logic from status monitoring
- Allows adding new monitoring features without modifying parking lot code
- Provides real-time updates to all interested parties
- Maintains single responsibility principle

### In Go, we can use models.ParkingLotObserver as a parameter type because it's an interface. This is powerful for several reasons:

1. Interface Acceptance: Any type that implements the ParkingLotObserver interface can be passed to this method. 
   - From the code, we know the interface looks like:
   ```
   type ParkingLotObserver interface {
       OnParkingLotStatusChanged(status ParkingLotStatus)
   }
   ```

2. Polymorphism: This means different types of observers can be used:
   - A MockObserver for testing
   - A ParkingManager for management
   - A Logger for logging
   - A DisplayBoard for UI updates
    As long as they implement OnParkingLotStatusChanged(), they can be observers.

3. Loose Coupling: The  ParkingLot doesn't need to know the concrete type of the observer. It only needs to know that the observer can handle notifications through OnParkingLotStatusChanged().
