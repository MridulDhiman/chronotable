# chronotable
in-memory hash table implementation with AOF persistence and versioning via point-in-time snapshotting supporting rollbacks and querying past states.

### Features
- In-memory storage using a Go map
- Thread safe operations via mutex locking.
- Client side persistence with AOF(Append only file)
- Snapshotting of the current state in binary file
- Versioning to support rollbacks and querying past states
- AOF Markers to each version for AOF based rollbacks

### Todo
- Upload snapshot and AOF files to object storage
- Scheduled backups of data
- Recovery in case of crash
- Diffing b/w diff. versions
- Ability to merge different versions
- branching as well
- checksums for verifying data integrity b/w the snapshots, using SHA256/1 digests
