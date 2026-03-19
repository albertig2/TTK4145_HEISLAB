# Elevator controller

The Elevator controller module contains all necessary routines and functions needed to control the hardware behavior of the local elevator. This includes functions for interactions with hardware inputs and outputs, routines for handling different events occurring during the lifetime of the elevator, and tools to manage and interact with the local order queue. The module consists of three sub modules, as well as the main finite state machine loop for the elevator: 

```` go
ElevatorController(ownId string, controllerChannels elevatorConfig.ElevatorControllerChannels, synchronisationChannels elevatorConfig.SynchronizationChannels, orderChannels elevatorConfig.OrderChannels)
````

## Elevator controller (Elevator controller.go)
The elevator controller contains one global function, the LocalElevatorController(). This is the finite state machine loop for the local elevator. The controller is built from the functions in the three sub modules, and ensures that the local elevator behaves correctly, and that the order assigner and synchronization module have the necessary information to serve their purpose. The controller mainly does this by detecting events and reacting accordingly by calling the correct handler. The controller has ten possible events that can occur, whereas some are local to the module and some are external and stemming from other modules:
-	Arrive at floor (local)
-	Receive order from button press (local)
-	Assigned order (external from order module)
-	Assigned order to peer (external from order module)
-	Order serviced by peer (external from order module)
-	Stop button event (local)
-	Obstruction event (local)
-	Timeout of open door timer(local)
-	Timeout for detect motor failure timer (local)
-	Timeout for send elevator update ticker (local)

The output and state changes are taken care of by the the handler function described below. The finite state machine in the controller has three possible states (called Behaviors): Idle, moving and open door.

### Functionality
- Implement the full controller module to create a plug and play elevator state machine. The state machine interacts with the surrounding code and modules through channels.

## Hardware (elevator_hardware.go)

The hardware sub module acts as an intermediate layer between the higher-level code and the hardware. It is essentially made to “hide” the raw hardware layer of the elevio hardware driver (elevio is pre made, and handed out as a project resource), and support the error handling in the finite state machine.

### Functionality
- Provide clean, clear and safe setters and getters for direct hardware interactions
- Provides some utility functions for executing several common and related hardware commands in series (like setting several order lights)
- Support for the main controller loop to detect motor failures

## Event handlers (elevator_controller_event_handlers.go)
The event handler contains functions to set the state machine output and decide the elevators’ reaction to the events that are generated either locally or externally. In this project events are defined loosely, and refer to inputs that warrants some sort of calculation or action from the elevator. The handlers are named after the event they correspond to. They will set the correct output based on the input, as well as change the state of the elevator in the finite state machine. This includes all reactions to button presses, handling the assigned orders received from the order assigner, ensuring correct light setting behavior and updating the order assigner of orders that has been serviced.
### Functionality
-	Routines for setting outputs based on inputs from channels, timers and the elevator object containing the internal states
-	Routines for initialization of the elevator object, the hardware and the controller channels
-	Support routines for the main event handlers 

## Order logic (local_elevator_order_logic.go)
The order logic sub module provide calculations needed to decide elevator behavior and direction based on local orders, as well as routines to manage the local order queue. 
### Functionality
-	Support for calculating the next motor direction and behavior based on the local order queue
-	Support for calculating which orders to clear based on behavior and local order queue
-	Functions for clearing orders in the local order queue based on previous calculations
-	Utility functions to detect orders relative to the elevator position
