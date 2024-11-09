# Chrono Table:  KV store with client side persistence to AOF with point-in-time snapshotting support

### Features:
- KV store which is storing data in hash table.
- Thread safe operations via mutex locking.
- Client side persistence with AOF(Append only file).
- Point in time snapshotting and versioning of the current state in binary file.
- Expiring a particular sample of keys  via cron job running every 1 sec.