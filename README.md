Msgstat
----------
Msgstat provides a CLI tool which analyses status-go log files and generates corresponding output of information 
related to messages delivered over whisper (either rpc or p2p).

## Commands

### `-readfile`

Msgstat `-readfile` expects to be provided a path to the proper log file which will be processed to generate reports.

```bash
> msgstat readfile ./logs/status-im/status-go/11.21.2017.log
```

### `-read`

Msgstat `-read` expects to be data from a log file written into `stdin` which will be processed to generate reports.

```bash
> cat ./logs/status-im/status-go/11.21.2017.log | msgstat read
```

## Report Format

Msgstat currently outputs parsed message into the following format:

```
Envelope #4T0NEM0Qrcm1uVVhBVmhYYW9HS1AwQ01xWmhwQmJGSDBnOHZpWnN2
FromDevice <4T0NEM0Qrcm1uVVhBVmhYYW9HS1A>
ToDevice <4T0NEM0Qrcm1uVVhBVmhYYW9HS1A>
Pipeline: RPC
Time: 2010-11-02
Status: Delivered
Payload: {
    Sell the wire around slack.
}

Envelope #4T0NEM0Qrcm1uVVhBVmhYYW9HS1AwQ01xWmhwQmJGSDBnOHZpWnN2
FromDevice <4T0NEM0Qrcm1uVVhBVmhYYW9HS1A>
ToDevice <4T0NEM0Qrcm1uVVhBVmhYYW9HS1A>
Pipeline: P2P
Time: 2010-11-02
Status: Delivered

Envelope #4T0NEM0Qrcm1uVVhBVmhYYW9HS1AwQ01xWmhwQmJGSDBnOHZpWnN2
FromDevice <4T0NEM0Qrcm1uVVhBVmhYYW9HS1A>
ToDevice <4T0NEM0Qrcm1uVVhBVmhYYW9HS1A>
Pipeline: P2P
Status: Pending
Time: 2010-11-02
Payload: {
    Back to freezia around slack.
}

Envelope #4T0NEM0Qrcm1uVVhBVmhYYW9HS1AwQ01xWmhwQmJGSDBnOHZpWnN2
FromDevice <4T0NEM0Qrcm1uVVhBVmhYYW9HS1A>
ToDevice <4T0NEM0Qrcm1uVVhBVmhYYW9HS1A>
Pipeline: P2P
Status: Failed
Time: 2010-11-02

```