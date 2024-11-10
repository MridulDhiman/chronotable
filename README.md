# chronotable
in-memory hash table implementation with AOF persistence and versioning via point-in-time snapshotting supporting rollbacks and querying past updates.

### Features
- In-memory storage using a Go map
- Thread safe operations via mutex locking.
- Client side persistence with AOF(Append only file)
- Snapshotting of the current state in binary file
- Versioning to support rollbacks and querying past states

### Todo
- Upload snapshot and AOF files to object storage
- Add AOF Markers to each version for AOF based rollbacks
- make `ChronoTable` a Singleton instance