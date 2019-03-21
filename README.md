# Errands Server
Errands API Server. A language agnostic, HTTP based queue. Persistant storage using Badger for SSD performance. Concurrency safe for multiple workers to be processing errands off the queue.

### Optional Params:
Errands server uses environment variables as a way to configure optional config params, here they are:
- `ERRANDS_PORT=:4545` - Will change the listening port to 4545
- `ERRANDS_STORAGE="/errands/"` - Will change the DB location to /errands/

### Running:
See the `postman` folder for the PostMan Collection which contains examples and tests for all the possible routes. 
