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
	// When defines the time when the notification was delivered.
	When time.Time `json:"when" toml:"when" yaml:"when"`
	// Status defines the current status of the element.
	Status string `json:"status" toml:"status" yaml:"status"`
	// Error defines the possible reason the message was rejected.
	Error string `json:"reason,omitempty" toml:"reason,omitempty" yaml:"error,omitempty"`
}

// AggregatedMessage contains aggregated facts for a given message envelope and it's total
// transition.
type AggregatedMessage struct {
	// Envelope holds the hash corresponding to the sent message.
	Envelope string `json:"envelope" toml:"envelope" yaml:"envelope"`
	// Protocol defines the means of transmission which is either P2P or RPC.
	Protocol string `json:"protocol" toml:"protocol" yaml:"protocol"`
	// FromDevice defines the public key hash which sent the message.
	FromDevice string `json:"from_device" toml:"from_device" yaml:"from_device"`
	// ToDevice defines the public key hash which will receive the message.
	ToDevice string `json:"to_device" toml:"to_device" yaml:"to_device"`
	// Payload defines the data attached to the message.
	Payload string `json:"payload" toml:"payload" yaml:"payload"`
	// Direction defines the route which a message comes through initially.
	Direction string `json:"direction" toml:"direction" yaml:"direction"`
	// Sent defines the time which the message was sent through the network.
	Sent time.Time `json:"sent_time" toml:"sent_time" yaml:"sent_time"`
	// Timeline defines the series of delivery changes that the message when through, this can be multiples in case of resending.
	Timeline []Timeline `json:"timeline" toml:"timeline" yaml:"timeline"`
	// Request defines the original whisper.NewMessage for this message, if found.
	Request whisper.NewMessage `json:"request" toml:"request" yaml:"request"`
}

// ReadAggregates processes incoming data from the reader and writes appropriate
// data with respect to format required.
func ReadAggregates(r io.ReadCloser, w io.WriteCloser, format string) error {
	defer w.Close()
	defer r.Close()

	format = strings.ToLower(format)

	aggregates, err := parseLogReader(r)
	if err != nil {
		return err
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

// parseLogReader parses out all log messages and returns associated AggregatedMessages
// or returns error if encountered.
func parseLogReader(r io.ReadCloser) ([]AggregatedMessage, error) {
	var order []string
	aggregates := make(map[string]AggregatedMessage)

	bufReader := bufio.NewReader(r)
	for {
		line, _, err := bufReader.ReadLine()

		// if we have reach end of file, just return.
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		logLine := string(line)
		if !messageHeader.MatchString(logLine) {
			continue
		}

		dataLog, err := decodeLogLine(logLine)
		if err != nil {
			continue
		}

		var message common.MessageState
		if err := json.Unmarshal([]byte(dataLog), &message); err != nil {
			continue
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

			order = append(order, message.Hash)
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

			msgAggr.Timeline = append(msgAggr.Timeline, Timeline{
				Status: "Pending",
				When:   message.Received,
			})
		case "Sent":
			if message.Type == outgoingMessage {
				msgAggr.Request = message.Source
			}

			msgAggr.Timeline = append(msgAggr.Timeline, Timeline{
				Status: "Sent",
				When:   message.Received,
			})
		case "Resent":
			msgAggr.Timeline = append(msgAggr.Timeline, Timeline{
				Status: "Resent",
				When:   message.Received,
			})
		case "Queued":
			msgAggr.Timeline = append(msgAggr.Timeline, Timeline{
				Status: "Queued",
				When:   message.Received,
			})
		case "Cached":
			msgAggr.Timeline = append(msgAggr.Timeline, Timeline{
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

			msgAggr.Timeline = append(msgAggr.Timeline, Timeline{
				Status: "Delivered",
				When:   message.Received,
			})
		case "Rejected":
			msgAggr.Timeline = append(msgAggr.Timeline, Timeline{
				Status: "Rejected",
				When:   message.Received,
				Error:  message.RejectionError,
			})
		case "Processing":
			msgAggr.Timeline = append(msgAggr.Timeline, Timeline{
				Status: "Processing",
				When:   message.Received,
			})
		}

		aggregates[message.Hash] = msgAggr
	}

	messages := make([]AggregatedMessage, len(order))
	for _, hash := range order {
		messages = append(messages, aggregates[hash])
	}

	return messages, nil
}

// decodeLogLine returns a decoded base64 status message which it strips off
// un-needed fields and uses the `state` field for it's decoding.
func decodeLogLine(line string) (string, error) {
	logLine := messageHeader.ReplaceAllLiteralString(line, "")
	logLine = strings.TrimSpace(logLine)
	logLine = strings.TrimPrefix(logLine, "geth=StatusIM")
	logLine = strings.TrimSpace(logLine)
	logLine = strings.TrimPrefix(logLine, "state=")

	if qlogLine, derr := strconv.Unquote(logLine); derr == nil {
		logLine = qlogLine
	}

	decoded, err := base64.StdEncoding.DecodeString(logLine)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
