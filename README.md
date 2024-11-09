# chronotable
An in-memory hash table with time travel capabilities. chronotable combines the speed of in-memory operations with the reliability of disk persistence and the power of temporal queries.

## Overview
chronotable is designed for applications that need fast key-value operations while maintaining historical state changes. It enables developers to query past states, roll back changes, and maintain data version history without sacrificing performance.

### Features
- In-memory storage using a Go map
- Thread safe operations via mutex locking.
- Client side persistence with AOF(Append only file)
- Snapshotting of the current state in binary file

### Current Targets
- Versioning to enable rollbacks and time-travel
- Upload snapshot and AOF files to object storage
