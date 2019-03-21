# Errands Server
Errands API Server. A language agnostic, HTTP based queue with server side events ( SSE ). Persistant storage using Badger for SSD performance. Concurrency safe for multiple workers to be processing errands off the queue.

## Client Libraries
- JS - https://github.com/polygon-io/errands-js
- Go - https://github.com/polygon-io/errands-go

## Optional Params:
Errands server uses environment variables as a way to configure optional config params, here they are:
- `ERRANDS_PORT=:4545` - Will change the listening port to 4545
- `ERRANDS_STORAGE="/errands/"` - Will change the DB location to /errands/

## Running:
See the `postman` folder for the PostMan Collection which contains examples and tests for all the possible routes. 

You can run the API locally with docker:

	docker run -p 5555:5555 polygonio/errands-server

## Push Events
There is an endpoint which will push server side events ( SSE ) to the client when changes happen on the errands server.

	/v1/errands/notifications?events={{events_to_listen_to}}

The events are comma delimited. For example `?events=created,completed,failed`. The default is `*` which will send any and all events which happen on the server.

Event Types:
- `created` - When a new errand is created
- `updated` - When a new errand is updated
- `completed` - When a new errand is completed
- `processing` - When a new errand is processing
- `retry` - When a new errand is being retried
- `failed` - When a new errand is failed

Along with the event type, the errand which triggered the event will have its data included in the message.