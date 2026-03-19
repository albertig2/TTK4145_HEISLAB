
stateDiagram-v2
    [*] --> NoOrder
    NoOrder: No order
    Pending: Pending
    Assigned: Assigned
    Complete: Complete

    NoOrder --> Pending: sees pending and not complete
    NoOrder --> Pending: sees on assigned
    NoOrder --> Complete: sees on complete
    Pending --> NoOrder: sees on complete
    Pending --> Assigned: sees on assigned
    Pending --> Pending: all pending
    Assigned --> Complete: completed
    Complete --> NoOrder: all other are here


<pre>
            +-----------+ 
+---------->| No order  |<--------------------------------+
|           +-----------+                                 |
| Other peer is   | Other peers are in Pending            |
| in Serviced     | and none are in Serviced              |
|                 v                                       |
|           +-----------+                                 |
+-----------|           |<---+                            |
            |  Pending  |    | Other peers are in Assigned|
            |           |----+                            |
            +-----------+                                 |
                  |                                       |
                  |                                       |
                  + All other peers are in Pending        |
                  |                                       |
                  |                                       |
                  v                                       |
            +-----------+                                 |
            | Assigned  |                                 |
            +-----------+                                 |
                  |                                       |
                  |                                       |
                  + Local peer arrived at floor           |
                  |                                       |
                  |                                       |
                  v                                       |
            +-----------+                                 |
            | Serviced  |---------------------------------+
            +-----------+ All other peers have No order
</pre>