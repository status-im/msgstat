package stats

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"

	whisper "github.com/ethereum/go-ethereum/whisper/whisperv5"
	"github.com/status-im/status-go/geth/common"
)

var (
	messageHeader   = regexp.MustCompile(`INFO \[\d+-\d+\|\d+:\d+:\d+\]\sMessage delivery notification`)
	incomingMessage = "IncomingMessage"
	outgoingMessage = "OutgoingMessage"
)

// Timeline contains timeline and status string for an aggregated message.
type Timeline struct {
	When   time.Time `json:"when" toml:"when" yaml:"when" csv:"when"`
	Status string    `json:"status" toml:"status" yaml:"status" csv:"status"`
	Error  string    `json:"reason,omitempty" toml:"reason,omitempty" yaml:"error,omitempty" csv:"error,omitempty"`
}

// AggregatedMessage contains aggregated facts for a given message envelope and it's total
// transition.
type AggregatedMessage struct {
	Envelope   string             `json:"envelope" toml:"envelope" yaml:"envelope" csv:"envelope"`
	Protocol   string             `json:"protocol" toml:"protocol" yaml:"protocol" csv:"protocol"`
	FromDevice string             `json:"from_device" toml:"from_device" yaml:"from_device" csv:"from_device"`
	ToDevice   string             `json:"to_device" toml:"to_device" yaml:"to_device" csv:"to_device"`
	Payload    string             `json:"payload" toml:"payload" yaml:"payload" csv:"payload"`
	Direction  string             `json:"direction" toml:"direction" yaml:"direction" csv:"direction"`
	Sent       time.Time          `json:"sent" toml:"sent" yaml:"sent" csv:"sent"`
	Timelines  []Timeline         `json:"timeline" toml:"timeline" yaml:"timeline" csv:"timeline"`
	Request    whisper.NewMessage `json:"request" toml:"request" yaml:"request" csv:"request"`
}

// ReadAggregate processes incoming data from the reader and writes appropriate
// data with respect to format required.
func ReadAggregate(r io.ReadCloser, w io.WriteCloser, format string) error {
	defer w.Close()
	defer r.Close()

	format = strings.ToLower(format)

	aggregates := make(map[string]AggregatedMessage)

	bufReader := bufio.NewReader(r)
	for {
		line, _, err := bufReader.ReadLine()

		// if we have reach end of file, just return.
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if !messageHeader.MatchString(string(line)) {
			continue
		}

		logLine := messageHeader.ReplaceAllLiteralString(string(line), "")
		logLine = strings.TrimSpace(logLine)
		logLine = strings.TrimPrefix(logLine, "geth=StatusIM")
		logLine = strings.TrimSpace(logLine)
		logLine = strings.TrimPrefix(logLine, "state=")

		if qlogLine, derr := strconv.Unquote(logLine); derr == nil {
			logLine = qlogLine
		}

		dataLog, err := base64.StdEncoding.DecodeString(logLine)
		if err != nil {
			return err
		}

		var message common.MessageState
		if err := json.Unmarshal([]byte(dataLog), &message); err != nil {
			return err
		}

		msgAggr, ok := aggregates[message.Hash]
		if !ok {
			msgAggr.Envelope = message.Hash
			msgAggr.Direction = message.Type
			msgAggr.Protocol = message.Protocol
			msgAggr.Sent = time.Unix(int64(message.TimeSent), 0)
			msgAggr.Payload = string(message.Payload)
			msgAggr.Request = message.Source
			msgAggr.FromDevice = message.FromDevice
			msgAggr.ToDevice = message.ToDevice
		}

		// This should be impossible, but skip any whoes Envelope hash matches but
		// has different protocols.
		if ok && msgAggr.Protocol != message.Protocol && msgAggr.Envelope != message.Hash {
			continue
		}

		if message.Type == outgoingMessage {
			if msgAggr.ToDevice == "" {
				msgAggr.ToDevice = message.ToDevice
			}

			if msgAggr.FromDevice == "" {
				msgAggr.FromDevice = message.FromDevice
			}

			if msgAggr.Payload == "" {
				msgAggr.Payload = string(message.Payload)
			}
		}

		switch message.Status {
		case "Pending":
			if message.Type == outgoingMessage {
				msgAggr.Request = message.Source
			}

			msgAggr.Timelines = append(msgAggr.Timelines, Timeline{
				Status: "Pending",
				When:   message.Received,
			})
		case "Sent":
			if message.Type == outgoingMessage {
				msgAggr.Request = message.Source
			}

			msgAggr.Timelines = append(msgAggr.Timelines, Timeline{
				Status: "Sent",
				When:   message.Received,
			})
		case "Resent":
			msgAggr.Timelines = append(msgAggr.Timelines, Timeline{
				Status: "Resent",
				When:   message.Received,
			})
		case "Queued":
			msgAggr.Timelines = append(msgAggr.Timelines, Timeline{
				Status: "Queued",
				When:   message.Received,
			})
		case "Cached":
			msgAggr.Timelines = append(msgAggr.Timelines, Timeline{
				Status: "Cached",
				When:   message.Received,
			})
		case "Delivered":
			if msgAggr.ToDevice == "" {
				msgAggr.ToDevice = message.ToDevice
			}

			if msgAggr.FromDevice == "" {
				msgAggr.FromDevice = message.FromDevice
			}

			msgAggr.Timelines = append(msgAggr.Timelines, Timeline{
				Status: "Delivered",
				When:   message.Received,
			})
		case "Rejected":
			msgAggr.Timelines = append(msgAggr.Timelines, Timeline{
				Status: "Rejected",
				When:   message.Received,
				Error:  message.RejectionError,
			})
		case "Processing":
			msgAggr.Timelines = append(msgAggr.Timelines, Timeline{
				Status: "Processing",
				When:   message.Received,
			})
		}

		aggregates[message.Hash] = msgAggr
	}

	for _, aggr := range aggregates {
		switch format {
		case "toml":
			if err := toml.NewEncoder(w).Encode(aggr); err != nil {
				return err
			}
		case "json":
			if err := json.NewEncoder(w).Encode(aggr); err != nil {
				return err
			}
		case "yaml":
			ymlData, err := yaml.Marshal(aggr)
			if err != nil {
				return err
			}

			w.Write(ymlData)
		}
	}

	return nil
}
