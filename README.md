# errands-server
Errands API Server. A language agnostic, HTTP based queue 

#### Formats:
Current response format:
```json
{
    "results": [
        {
            "id": "718f112b-56d8-48bf-b626-681eb0b7ed7a",
            "name": "New Errand Job",
            "type": "backfill",
            "options": {
                "ttl": 60
            },
            "created": 1552936821129,
            "status": "inactive"
        },
        {
            "id": "d9bde463-1a77-4356-bd9a-1eb001335e6c",
            "name": "Process: Banner Ads v1",
            "type": "sv-extract",
            "options": {
                "ttl": 60
            },
            "data": {
                "dbid": "899687969",
                "file": "http://s3.stackdot.com/some/dir/file.sketch",
                "user": "351351345151"
            },
            "created": 1552937527388,
            "status": "inactive"
        }
    ],
    "status": "OK"
}
```