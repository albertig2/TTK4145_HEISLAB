# Order Protocol
This document describes the distributed order protocol for a multi-elevator system. The protocol ensures robust and fault-tolerant assignment and completion of hall and cab orders, even in the presence of network failures or elevator crashes.

## State Description
Unknown state is only for cab orders and Serviced is only for hall orders. The other states are common for both cab and hall orders. 
- Unknown: After initializing or reinitializing, cab orders are in an unknown state until confirmation is received from other elevators or until a new order is registered. If the elevator is alone, it remains in the unknown state until it receives an order.
- NoOrder: No active order for this floor/button.
- Pending: Order is known by the elevator but not yet assigned. The order can be assigned when all elvators are in Pending.
- Assigned: The order is assigned to the elevator with the order in Assigned state.
- Serviced: The order is serviced by the elevator with the order in Serviced state.

The hall orders are initialized to NoOrder and the cab orders go straight from Assigned to NoOrder when the order is serviced.

## Hall Orders
Hall orders are generated when a user presses a hall button (up/down) on any floor. These orders are shared across all elevators in the distributed system. When a hall order is registered, every online elevator must acknowledge it by setting the corresponding order to Pending. Once all elevators have acknowledged the order (i.e., all are in Pending), the system assigns the order to one elevator, which transitions the order to Assigned. The assignment only considers elevators that are up and running. If the assigned elevator fails or disconnects, all remaining elevators will have the order in Pending, and the order will be reassigned. This mechanism ensures that hall orders are always eventually serviced, even in the presence of network failures or elevator crashes.

### Hall Transitions 
- NoOrder --> Pending: Other peers in Pending or Assigned and none in Serviced or local hall button is pressed
- Pending --> Pending: Other peers are in Assigned
- Pending --> Assigned: All other peers are in Pending
- Assigned --> Serviced: Local peer arrived at floor
- Serviced --> NoOrder: All other peers are in NoOrder
- Pending --> NoOrder: Other peers are in Serviced

<pre>
                        +-----------+ 
            +---------->| No order  |<--------------------------------+
            |           +-----------+                                 |
            | Other peer is   | Other peers are in Pending or Assigned|
            | in Serviced     |      and none are in Serviced         |
            |                 |   or local hall button is pressed     |
            |                 v                                       |
            |           +-----------+                                 |
            +-----------|           |----+                            |
                        |  Pending  |    | Other peers are in Assigned|
                        |           |<---+                            |
                        +-----------+                                 |
                              |                                       |
                              |                                       |
                              | All other peers are in Pending        |
                              |                                       |
                              |                                       |
                              v                                       |
                        +-----------+                                 |
                        | Assigned  |                                 |
                        +-----------+                                 |
                              |                                       |
                              |                                       |
                              | Local peer arrived at floor           |
                              |                                       |
                              |                                       |
                              v                                       |
                        +-----------+                                 |
                        | Serviced  |---------------------------------+
                        +-----------+ All other peers have No order
</pre>

## Cab Orders
Cab orders are generated when a user presses a cab button inside an elevator. These orders are initially only known to the local elevator, but are also shared on the network with other elevators for redundancy. If an elevator is restarted or rejoins the network, the other elevators will resend any outstanding cab orders to ensure no requests are lost. This protocol guarantees that all cab orders are eventually serviced, even if elevators are restarted or temporarily disconnected. 

### Cab Transitions
- Unknown --> NoOrder: All other peers are in NoOrder
- Unknown --> Pending: Other peers are in Pending and none in Assigned
- Unknown --> Assigned: All other peers are in Pending or Assigned
- NoOrder --> Pending: Local cab button is pressed
- Pending --> Assigned: All other peers are in Pending
- Assigned --> NoOrder: Local peer arrived at floor


<pre>
                  +-----------+ 
                  | Unknown   |------------------------------+
                  +-----------+                              |
                        |                                    |
                        | All other peers are in NoOrder     | 
                        |                                    |
                        v                                    |  
                  +-----------+                              |
+---------------->|  NoOrder  |                              |
|                 +-----------+                              | Other peers are in Pending
|                       |                                    |   and none in Assigned
|                       | Local cab button is pressed        |
|                       |                                    |
| Local peer            v                                    |
|  arrived        +-----------+                              |
|                 | Pending   |<-----------------------------+
|                 +-----------+                              |
|                       |                                    |
|                       | All other peers are in Pending     | All other peers are in
|                       |                                    |   Pending or Assigned
|                       v                                    |
|                 +-----------+                              |
+-----------------| Assigned  |<-----------------------------+
                  +-----------+
</pre>