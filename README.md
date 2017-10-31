Msgstat
----------
Msgstat provides a CLI tool which analyses status-go log files and generates corresponding output of information
related to messages delivered over whisper (either rpc or p2p).

## Commands

### `readfile`

Msgstat `readfile` expects to be provided a path to the proper log file which will be processed to generate reports.

```bash
> msgstat readfile ./logs/status-im/status-go/11.21.2017.log
```

### `read`

Msgstat `read` expects to be data from a log file written into `stdin` which will be processed to generate reports.

```bash
> cat ./logs/status-im/status-go/11.21.2017.log | msgstat read
```

## Report Format

Msgstat currently outputs parsed message into the following format:

```json
{
  "envelope_hash": "#4T0NEM0Qrcm1uVVhBVmhYYW9HS1AwQ01xWmhwQmJGSDBnOHZpWnN2",
  "from_device": "4T0NEM0Qrcm1uVVhBVmhYYW9HS1A",
  "to": "4T0NEM0Qrcm1uVVhBVmhYYW9HS1A",
  "protocol": "rpc",
  "time": 2010-11-12,
  "status-diffs": [
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

{
  "envelope_hash": "#4T0NEM0Qrcm1uVVhBVmhYYW9HS1AwQ01xWmhwQmJGSDBnOHZpWnN2",
  "from_device": "4T0NEM0Qrcm1uVVhBVmhYYW9HS1A",
  "to": "4T0NEM0Qrcm1uVVhBVmhYYW9HS1A",
  "protocol": "rpc",
  "time": 2010-11-12,
  "status-diffs": [
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
