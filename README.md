Generic map wrapper for client side persistence via AOF and point in time snapshotting

### Features:
- Generic KV store for storing data in hash table.
- Client side persistence with AOF(Append only file).
- Point in time snapshotting and versioning of the current state in binary file.
- Expiring a particular sample of keys  via cron job running every 1 sec.