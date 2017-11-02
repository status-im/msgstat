Msgstat
----------
Msgstat provides a CLI tool which analyses status-go log files and generates corresponding aggregated data
related to messages delivery over whisper (either rpc or p2p) protocols.

## Usage

Msgstat spits out structures that provide commulation of all related facts about a given message based on desired output format supported.

*It by defaults expects processed data to be delivered either through pipes through `stdin` or through the use of the `file` flag.*

- Processing log lines through stdin

```bash
> cat ./logs/status-im/status-go/11.21.2017.log | msgstat
```

- Processing log lines through stdin with `yaml` output format

```bash
> cat ./logs/status-im/status-go/11.21.2017.log | msgstat -format=yaml
```

- Processing log lines from a file

```bash
> msgstat -file=./logs/status-im/status-go/11.21.2017.log
```

All processed data is written through `stdout`, except when the use of the `-out` flag is used to point to a output file to contain
created data.

### Flags

#### `-file=./status-log`

This flag is used to indicate that file should be used to retrieve all log lines to be processed.

#### `-out=./processed-log.json`

This flag is used to indicate that all processed data should be stored into provided path.

#### `-format=json|yaml|toml`

This flag is used to dictate the encoding format which should be used for aggregated data.


### Log Input Format

Msgstat expects that all log contain log lines separated by new lines, and targets log lines such has below:

```
INFO [10-31|14:57:13] Message delivery notification            geth=StatusIM state=eyJ0eXBlIjoiSW5jb21pbmdNZXNzYWdlIiwicHJvdG9jb2wiOiJSUEMiLCJzdGF0dXMiOiJDYWNoZWQiLCJlbnZlbG9wZSI6IkJCcHNCQTFiZ0RXNXUyQ0lWL3hyRlcvaThFUG5mZGxwdmxCcnNrbnUwN3ZZVXBhNXRnT3NVWUdhbmljYmY1MjdpOEg5UmRTWERtWUhOSHlyakVSOS9pQkgrcGxWS3FmME9mMHlRMlRKSHB2ZklrUU9VV1g0dzdENE1mTmZXYVUycER3djJQUG43VjVjaDRJeE9Oajh4YnFwMjhMQnErMzhmczBxaG9CYzVLWTczcnZ0T1NsNm12dFRWcGg2RENsUEVQOE1zU2pvaDBYakdXL3lYTnpBeWxUODNnZzdtdFRrY1l3V0ZhSHk5OEV5V1lVTFg5WHBTeUJFVVVVQ25GWUVoRDV1ZzBaUkZyY2ltN29QTng2RWZVZUhuc3lXLysvZFFRb1FMcXI4K1lvRFJaUVZSTVVDSUZwbWFBRWVjQUkvVkFSMGhsTkdrV2pGY0JkNzVjZXdIZzJvWWtuKzgrQlA3bTJsL0FKdk1COFJMVW5tZlpBN1dlOUJyUnM5alFKZGJiSmpCcGRaYVNqS0ZoVCtTYnhvVi9oWllobmFFVDQyZzZJajh3d1hGRHNNck54MlpJUEcxMEZhdm02Y0JlNWgybnh2ME55ZHVmaS92eFZaaTRuRU9SZGpuQ0RRZlhab0Rab1EwSXpvS1N5VyIsInRpbWUiOjE1MDk0NTgyMzMsImVudmVsb3BlX2hhc2giOiIweDU1YmMyMmJlNGIxYTQxOTRhZmQ4OTNlZTYzZDA1ZjU4MmIwYzAxMTk3OGJkMzkxM2FlYTdlYmI4Mjc4ZmMwZGIiLCJzb3VyY2UiOnsic3ltS2V5SUQiOiIiLCJwdWJLZXkiOiIweCIsInNpZyI6IiIsInR0bCI6MCwidG9waWMiOiIweDAwMDAwMDAwIiwicGF5bG9hZCI6IjB4IiwicGFkZGluZyI6IjB4IiwicG93VGltZSI6MCwicG93VGFyZ2V0IjowLCJ0YXJnZXRQZWVyIjoiIn19
```

Msgstat will extract any log lines with matching `Message Delivery Notification` header, extracting the `state` field and it's value which contains base64 encoding delivery data.


### Output Format

Msgstat currently supports the following formats:

- `TOML`

```toml
envelope = "0x156789a7892cfa5b5a45ca9f3187799d2a0293034175093ab8726ecbd8cbc6c6"
protocol = "RPC"
from_device = "0x04eedbaafd6adf4a9233a13e7b1c3c14461fffeba2e9054b8d456ce5f6ebeafadcbf3dce3716253fbc391277fa5a086b60b283daf61fb5b1f26895f456c2f31ae3"
to_device = ""
payload = ""
direction = "OutgoingMessage"
sent_time = 1970-01-01T00:00:00Z

[[timeline]]
  when = 2017-11-01T14:35:28Z
  status = "Pending"

[[timeline]]
  when = 2017-11-01T14:35:32Z
  status = "Pending"

[[timeline]]
  when = 2017-11-01T14:35:33Z
  status = "Pending"

[[timeline]]
  when = 2017-11-01T14:35:34Z
  status = "Pending"

[request]
  SymKeyID = "d29dd5c6470d556c20e69b0f01b269b00c74e363cae442628171256a6df5ce38"
  PublicKey = []
  Sig = ""
  TTL = 20
  Topic = "0x0f1a4771"
  Payload = [116, 101, 115, 116, 32, 109, 101, 115, 115, 97, 103, 101, 32, 52, 32, 40, 34, 34, 32, 45, 62, 32, 34, 34, 44, 32, 97, 110, 111, 110, 32, 98, 114, 111, 97, 100, 99, 97, 115, 116, 41]
  Padding = []
  PowTime = 20
  PowTarget = 0.01
  TargetPeer = ""
```

- `YAML`

```yaml
envelope: 0x156789a7892cfa5b5a45ca9f3187799d2a0293034175093ab8726ecbd8cbc6c6
protocol: RPC
from_device: 0x04eedbaafd6adf4a9233a13e7b1c3c14461fffeba2e9054b8d456ce5f6ebeafadcbf3dce3716253fbc391277fa5a086b60b283daf61fb5b1f26895f456c2f31ae3
to_device: ""
payload: ""
direction: OutgoingMessage
sent_time: 1970-01-01T01:00:00+01:00
timeline:
- when: 2017-11-01T15:35:28.139031412+01:00
  status: Pending
- when: 2017-11-01T15:35:32.226407285+01:00
  status: Pending
- when: 2017-11-01T15:35:33.357611609+01:00
  status: Pending
- when: 2017-11-01T15:35:34.470666086+01:00
  status: Pending
request:
  symkeyid: d29dd5c6470d556c20e69b0f01b269b00c74e363cae442628171256a6df5ce38
  publickey: []
  sig: ""
  ttl: 20
  topic: "0x0f1a4771"
  payload:
  - 116
  - 101
  - 115
  - 116
  - 32
  - 109
  - 101
  - 115
  - 115
  - 97
  - 103
  - 101
  - 32
  - 52
  - 32
  - 40
  - 34
  - 34
  - 32
  - 45
  - 62
  - 32
  - 34
  - 34
  - 44
  - 32
  - 97
  - 110
  - 111
  - 110
  - 32
  - 98
  - 114
  - 111
  - 97
  - 100
  - 99
  - 97
  - 115
  - 116
  - 41
  padding: []
  powtime: 20
  powtarget: 0.01
  targetpeer: ""
```

- `JSON`

```json
{
  "envelope": "0xa1182bd04fd5d60717b3ae0d19569a65b00340f2fb1ba31b99873eb4b34a66f0",
  "protocol": "RPC",
  "from_device": "",
  "to_device": "",
  "payload": "",
  "direction": "OutgoingMessage",
  "sent_time": "2017-11-01T15:35:34+01:00",
  "timeline": [
    {
      "when": "2017-11-01T15:35:34.470932867+01:00",
      "status": "Sent"
    },
    {
      "when": "2017-11-01T15:35:34.47104423+01:00",
      "status": "Queued"
    },
    {
      "when": "2017-11-01T15:35:34.480469935+01:00",
      "status": "Rejected",
      "reason": "processing message: does not match"
    },
    {
      "when": "2017-11-01T15:35:34.471033731+01:00",
      "status": "Cached"
    },
    {
      "when": "2017-11-01T15:35:34.48042816+01:00",
      "status": "Processing"
    },
    {
      "when": "2017-11-01T15:35:34.480454857+01:00",
      "status": "Rejected",
      "reason": "Envelope failed to be opened"
    }
  ],
  "request": {
    "symKeyID": "d29dd5c6470d556c20e69b0f01b269b00c74e363cae442628171256a6df5ce38",
    "pubKey": "0x",
    "sig": "",
    "ttl": 20,
    "topic": "0x0f1a4771",
    "payload": "0x74657374206d657373616765203420282222202d3e2022222c20616e6f6e2062726f61646361737429",
    "padding": "0x",
    "powTime": 20,
    "powTarget": 0.01,
    "targetPeer": ""
  }
}
```
