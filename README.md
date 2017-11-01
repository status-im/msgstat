Msgstat
----------
Msgstat provides a CLI tool which analyses status-go log files and generates corresponding output of information
related to messages delivered over whisper (either rpc or p2p).

## Usage

Msgstat provides the processing functionality for the processing of all log files which would be used to produced aggregated collection of message delivery facts.

It by defaults expects processing data to be delivered either through sending given logs through `stdin` or through the use of the `file` flag.

- Processing log lines through stdin

```bash
> cat ./logs/status-im/status-go/11.21.2017.log | msgstat
```

- Processing log lines from a file

```bash
> msgstat read -file=./logs/status-im/status-go/11.21.2017.log
```

Where all information is delivered through the `stdout`, except when the use of the `-out` flag is used.

### Flags

#### `-file=./status-log`
This flag is used to indicate that all a file should be used to retrieve all log lines to be processed line by line for the generation of processed log facts.

#### `-out=./processed-log.json`
This flag is used to indicate that all processed data should be stored into provided file.

#### `-format=json|yaml|csv`

This flag is used to dictate the encoding format which should be used for generation of processed resource which would
then be processed by external tools as needed.


### Log Format

Msgstat expects that all log contain log lines seperated by new lines, and it targets logs which mark this format:

```
INFO [10-31|14:57:13] Message delivery notification            geth=StatusIM state=eyJ0eXBlIjoiSW5jb21pbmdNZXNzYWdlIiwicHJvdG9jb2wiOiJSUEMiLCJzdGF0dXMiOiJDYWNoZWQiLCJlbnZlbG9wZSI6IkJCcHNCQTFiZ0RXNXUyQ0lWL3hyRlcvaThFUG5mZGxwdmxCcnNrbnUwN3ZZVXBhNXRnT3NVWUdhbmljYmY1MjdpOEg5UmRTWERtWUhOSHlyakVSOS9pQkgrcGxWS3FmME9mMHlRMlRKSHB2ZklrUU9VV1g0dzdENE1mTmZXYVUycER3djJQUG43VjVjaDRJeE9Oajh4YnFwMjhMQnErMzhmczBxaG9CYzVLWTczcnZ0T1NsNm12dFRWcGg2RENsUEVQOE1zU2pvaDBYakdXL3lYTnpBeWxUODNnZzdtdFRrY1l3V0ZhSHk5OEV5V1lVTFg5WHBTeUJFVVVVQ25GWUVoRDV1ZzBaUkZyY2ltN29QTng2RWZVZUhuc3lXLysvZFFRb1FMcXI4K1lvRFJaUVZSTVVDSUZwbWFBRWVjQUkvVkFSMGhsTkdrV2pGY0JkNzVjZXdIZzJvWWtuKzgrQlA3bTJsL0FKdk1COFJMVW5tZlpBN1dlOUJyUnM5alFKZGJiSmpCcGRaYVNqS0ZoVCtTYnhvVi9oWllobmFFVDQyZzZJajh3d1hGRHNNck54MlpJUEcxMEZhdm02Y0JlNWgybnh2ME55ZHVmaS92eFZaaTRuRU9SZGpuQ0RRZlhab0Rab1EwSXpvS1N5VyIsInRpbWUiOjE1MDk0NTgyMzMsImVudmVsb3BlX2hhc2giOiIweDU1YmMyMmJlNGIxYTQxOTRhZmQ4OTNlZTYzZDA1ZjU4MmIwYzAxMTk3OGJkMzkxM2FlYTdlYmI4Mjc4ZmMwZGIiLCJzb3VyY2UiOnsic3ltS2V5SUQiOiIiLCJwdWJLZXkiOiIweCIsInNpZyI6IiIsInR0bCI6MCwidG9waWMiOiIweDAwMDAwMDAwIiwicGF5bG9hZCI6IjB4IiwicGFkZGluZyI6IjB4IiwicG93VGltZSI6MCwicG93VGFyZ2V0IjowLCJ0YXJnZXRQZWVyIjoiIn19
```

Msgstat will extract any log lines with matching `Mesage Delivery Notification` and extract the `state` field and value which contains the needed delivery facts.


### Report Format

Msgstat currently outputs parsed message into the following format:

```json
{
  "protocol": "RPC",
  "envelope": "#4T0NEM0Qrcm1uVVhBVmhYYW9HS1AwQ01xWmhwQmJGSDBnOHZpWnN2",
  "from_device": "4T0NEM0Qrcm1uVVhBVmhYYW9HS1A",
  "to_device": "4T0NEM0Qrcm1uVVhBVmhYYW9HS1A",
  "time": 2010-11-12,
  "timelines": [
    {
      "time_stamp": "2010-11-12 10:20:100T34343",
      "status": "pending"
    },
    {
      "time_stamp": "2010-11-12 3:20:100T34343",
      "status": "rejected",
      "reason": "Failed to locate peer"
    }
  ],
  "payload": "Thunder crash",
}
```
