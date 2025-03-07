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

## Strategy Pattern

### What is the Strategy Pattern?

The **Strategy pattern** is a behavioral design pattern that allows you to define a family of algorithms, put each of them into separate classes, and make them interchangeable. It lets the algorithm vary independently from the clients that use it. Essentially, it’s a way to encapsulate different ways of doing something and switch between them at runtime.

Think of it like choosing a tool for a job: you might have a screwdriver, a hammer, or a wrench, and you pick the one that best suits the task at hand. The Strategy pattern formalizes this idea in code.

---

### Key Components

1. **Strategy (Interface or Abstract Class)**:
   - Defines a common interface that all concrete strategies must follow.
   - Typically has a method that encapsulates the algorithm/behavior.

2. **Concrete Strategies**:
   - The actual implementations of the algorithm.
   - Each class implements the Strategy interface differently.

3. **Context**:
   - The class that uses the strategy.
   - Holds a reference to a Strategy object and can switch between strategies.

---

### How It Works

- The Context delegates the work to the Strategy object instead of implementing the behavior itself.
- You can change the strategy (algorithm) at runtime by passing a different Concrete Strategy to the Context.
- The client (code using the Context) doesn’t need to know the details of how the strategy works—it just knows it gets the job done.

---

### Real-World Analogy

Imagine you’re traveling from point A to point B. You could:
- Take a car (fast but costly).
- Walk (free but slow).
- Ride a bike (moderate speed and effort).

The **goal** is to get to your destination (the Context), and the **strategy** is how you choose to travel (car, walking, bike). You can switch your travel method based on weather, time, or preference without changing your overall plan.

---

### When to Use It

- When you have multiple ways to perform a task and want to avoid lots of conditional statements (like `if-else` or `switch`).
- When you want to isolate algorithm-specific code from the rest of your application.
- When you need to switch between algorithms at runtime.
- When you want to make your code more extensible (easy to add new strategies).

---

### Benefits

1. **Flexibility**: Easily swap or add new algorithms without changing the core code.
2. **Single Responsibility**: Each strategy handles one specific way of doing things.
3. **Open/Closed Principle**: You can extend behavior (add new strategies) without modifying existing code.
4. **Cleaner Code**: Avoids messy conditional logic by encapsulating algorithms.

---

### Drawbacks

1. **Increased Complexity**: Adds more classes/interfaces, which can overcomplicate simple scenarios.
2. **Client Awareness**: The client might need to know about the available strategies to choose one.

---

### Example in the Go Code I Shared Earlier

In the payment processing example:
- **Strategy Interface**: `PaymentStrategy` with the `Pay()` method.
- **Concrete Strategies**: `CreditCardPayment`, `PayPalPayment`, `CryptoPayment`—each implements `Pay()` differently.
- **Context**: `ShoppingCart`, which uses a `PaymentStrategy` to process payments and can switch between them with `SetPaymentStrategy()`.

You can pay with a credit card today and switch to PayPal tomorrow without changing how the `ShoppingCart` works—it just delegates the payment logic to the chosen strategy.

---

### Comparison to Other Patterns

- **Strategy vs. State**: Strategy focuses on interchangeable algorithms, while State focuses on changing behavior based on an object’s internal state.
- **Strategy vs. Template Method**: Strategy uses composition (delegating to a strategy object), while Template Method uses inheritance to define an algorithm’s skeleton.
