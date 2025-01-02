# chronotable
in-memory hash table implementation with AOF persistence and versioning via point-in-time snapshotting supporting rollbacks and querying past states.

### Features
- In-memory storage using a Go map
- Thread safe operations via mutex locking.
- write ahead logging of put/delete operations to Append Only File(AOF).
- Binary serialization/deserialization of current snapshot state.
- Versioning to support rollbacks and querying past states.
- Replay writes in case of crash
- Start from latest state in case of restart
- Viper based local and remote configuration management 

### Todo
- Scheduled backups of data.
- Diffing b/w diff. versions.
- Upload snapshots and AOF files to object storage/NoSQL database.
- Ability to merge different versions.
- support branching.
- checksums for verifying data integrity b/w the snapshots using SHA256/1 digests.
- writes blackbox/whitebox unit tests for each package. 