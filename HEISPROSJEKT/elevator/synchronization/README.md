# Synchronization Module

## Overview
This module is responsible for synchronizing elevator state information across multiple peers in a distributed elevator system. It maintains a shared representation of all elevators (`PeerView`) and ensures consistent updates through message passing and periodic broadcasting.

The module integrates local elevator updates, peer updates, and network communication to maintain a coherent system state.

## Responsibilities
- Maintain a local representation of all known peers and their states
- Integrate updates from:
  - Local elevator
  - Other peers
  - Network broadcasts
- Distribute the local system state to other peers
- Track currently alive peers in the system

## File Structure

### `elevator_system.go`
Handles manipulation and maintenance of the shared system state (`PeerView`).
- Provides setter functions for updating local elevator state
- Initializes and maintains peer state structures
- Merges incoming peer information into the local view

### `peer_synchronization_functions.go`
Defines helper functions related to synchronization and channel setup.
- Initializes all synchronization-related channels
- Updates the local system state based on elevator input

### `synchronize_elevators.go`
Implements the main synchronization loop.
- Coordinates communication between modules using channels
- Handles incoming peer messages and system updates
- Periodically broadcasts the local `PeerView`