# chronotable
in-memory hash table implementation with AOF persistence and versioning via point-in-time snapshotting supporting rollbacks and querying past states.

### Features
- In-memory storage using a Go map
- Thread safe operations via mutex locking.
- Client side persistence with AOF(Append only file)
- Binary serialization/deserialization of current snapshot state
- Versioning to support rollbacks and querying past states
- AOF Markers to each version for AOF based rollbacks
- Replay writes in case of crash
- Start from latest state in case of restart
- View current state changes
- Viper based configuration management 

### Todo
- Buffered Channel/Wait Group which will block the main thread till all the goroutines are not completed
- Implement transactions for each commit using AOF file
- Upload snapshot and AOF files to object storage/NoSQL database
- Scheduled backups of data
- Diffing b/w diff. versions
- Ability to merge different versions
- branching as well
- checksums for verifying data integrity b/w the snapshots, using SHA256/1 digests
