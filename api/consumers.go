// Copyright 2020 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"encoding/json"
	"fmt"
	"time"
)

// also update wellKnownSubjectSchemas
const (
	JSApiConsumerCreateT                   = "$JS.API.CONSUMER.CREATE.%s"
	JSApiConsumerCreate                    = "$JS.API.CONSUMER.CREATE.*"
	JSApiConsumerCreateWithNameT           = "$JS.API.CONSUMER.CREATE.%s.%s"
	JSApiConsumerCreateWithName            = "$JS.API.CONSUMER.CREATE.*.>"
	JSApiDurableCreateT                    = "$JS.API.CONSUMER.DURABLE.CREATE.%s.%s"
	JSApiDurableCreate                     = "$JS.API.CONSUMER.DURABLE.CREATE.*.*"
	JSApiConsumerCreateExT                 = "$JS.API.CONSUMER.CREATE.%s.%s.%s"
	JSApiConsumerCreateEx                  = "$JS.API.CONSUMER.CREATE.%s.*.>"
	JSApiConsumerNamesT                    = "$JS.API.CONSUMER.NAMES.%s"
	JSApiConsumerNames                     = "$JS.API.CONSUMER.NAMES.*"
	JSApiConsumerListT                     = "$JS.API.CONSUMER.LIST.%s"
	JSApiConsumerList                      = "$JS.API.CONSUMER.LIST.*"
	JSApiConsumerInfoT                     = "$JS.API.CONSUMER.INFO.%s.%s"
	JSApiConsumerDeleteT                   = "$JS.API.CONSUMER.DELETE.%s.%s"
	JSApiRequestNextT                      = "$JS.API.CONSUMER.MSG.NEXT.%s.%s"
	JSApiRequestNext                       = "$JS.API.CONSUMER.MSG.NEXT.*.*"
	JSApiConsumerLeaderStepDownT           = "$JS.API.CONSUMER.LEADER.STEPDOWN.%s.%s"
	JSApiConsumerPauseT                    = "$JS.API.CONSUMER.PAUSE.%s.%s"
	JSApiConsumerUnpinT                    = "$JS.API.CONSUMER.UNPIN.%s.%s"
	JSMetricConsumerAckPre                 = JSMetricPrefix + ".CONSUMER.ACK"
	JSAdvisoryConsumerMaxDeliveryExceedPre = JSAdvisoryPrefix + ".CONSUMER.MAX_DELIVERIES"
)

type ConsumerAction int

const (
	ActionCreateOrUpdate ConsumerAction = iota
	ActionUpdate
	ActionCreate
)

const (
	actionUpdateString         = "update"
	actionCreateString         = "create"
	actionCreateOrUpdateString = ""
)

func (a ConsumerAction) String() string {
	switch a {
	case ActionCreateOrUpdate:
		return actionCreateOrUpdateString
	case ActionCreate:
		return actionCreateString
	case ActionUpdate:
		return actionUpdateString
	}
	return actionCreateOrUpdateString
}
func (a ConsumerAction) MarshalJSON() ([]byte, error) {
	switch a {
	case ActionCreate:
		return json.Marshal(actionCreateString)
	case ActionUpdate:
		return json.Marshal(actionUpdateString)
	case ActionCreateOrUpdate:
		return json.Marshal(actionCreateOrUpdateString)
	default:
		return nil, fmt.Errorf("can not marshal %v", a)
	}
}

func (a *ConsumerAction) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case jsonString("create"):
		*a = ActionCreate
	case jsonString("update"):
		*a = ActionUpdate
	case jsonString(""):
		*a = ActionCreateOrUpdate
	default:
		return fmt.Errorf("unknown consumer action: %v", string(data))
	}
	return nil
}

// io.nats.jetstream.api.v1.consumer_unpin_response
type JSApiConsumerUnpinRequest struct {
	Group string `json:"group"`
}

// io.nats.jetstream.api.v1.consumer_unpin_response
type JSApiConsumerUnpinResponse struct {
	JSApiResponse
}

// io.nats.jetstream.api.v1.consumer_delete_response
type JSApiConsumerDeleteResponse struct {
	JSApiResponse
	Success bool `json:"success,omitempty"`
}

// io.nats.jetstream.api.v1.consumer_create_request
type JSApiConsumerCreateRequest struct {
	Stream string         `json:"stream_name"`
	Config ConsumerConfig `json:"config"`
	Action ConsumerAction `json:"action"`
	// Pedantic disables server features that would set defaults and adjust the provided config
	Pedantic bool `json:"pedantic,omitempty"`
}

// io.nats.jetstream.api.v1.consumer_create_response
type JSApiConsumerCreateResponse struct {
	JSApiResponse
	*ConsumerInfo
}

// io.nats.jetstream.api.v1.consumer_info_response
type JSApiConsumerInfoResponse struct {
	JSApiResponse
	*ConsumerInfo
}

// io.nats.jetstream.api.v1.consumer_names_request
type JSApiConsumerNamesRequest struct {
	JSApiIterableRequest
}

// io.nats.jetstream.api.v1.consumer_names_response
type JSApiConsumerNamesResponse struct {
	JSApiResponse
	JSApiIterableResponse
	Consumers []string `json:"consumers"`
}

// io.nats.jetstream.api.v1.consumer_list_request
type JSApiConsumerListRequest struct {
	JSApiIterableRequest
}

// io.nats.jetstream.api.v1.consumer_list_response
type JSApiConsumerListResponse struct {
	JSApiResponse
	JSApiIterableResponse
	Consumers []*ConsumerInfo `json:"consumers"`
	Missing   []string        `json:"missing,omitempty"`
}

// io.nats.jetstream.api.v1.consumer_leader_stepdown_request
type JSApiConsumerLeaderStepdownRequest struct {
	Placement *Placement `json:"placement,omitempty"`
}

// io.nats.jetstream.api.v1.consumer_leader_stepdown_response
type JSApiConsumerLeaderStepDownResponse struct {
	JSApiResponse
	Success bool `json:"success,omitempty"`
}

// io.nats.jetstream.api.v1.consumer_pause_request
type JSApiConsumerPauseRequest struct {
	PauseUntil time.Time `json:"pause_until,omitempty"`
}

// io.nats.jetstream.api.v1.consumer_pause_response
type JSApiConsumerPauseResponse struct {
	JSApiResponse
	Paused         bool          `json:"paused"`
	PauseUntil     time.Time     `json:"pause_until"`
	PauseRemaining time.Duration `json:"pause_remaining,omitempty"`
}

type AckPolicy int

const (
	AckNone AckPolicy = iota
	AckAll
	AckExplicit
)

func (p AckPolicy) String() string {
	switch p {
	case AckNone:
		return "None"
	case AckAll:
		return "All"
	case AckExplicit:
		return "Explicit"
	default:
		return "Unknown Acknowledgement Policy"
	}
}

func (p *AckPolicy) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case jsonString("none"):
		*p = AckNone
	case jsonString("all"):
		*p = AckAll
	case jsonString("explicit"):
		*p = AckExplicit
	default:
		return fmt.Errorf("can not unmarshal %q", data)
	}

	return nil
}

func (p AckPolicy) MarshalJSON() ([]byte, error) {
	switch p {
	case AckNone:
		return json.Marshal("none")
	case AckAll:
		return json.Marshal("all")
	case AckExplicit:
		return json.Marshal("explicit")
	default:
		return nil, fmt.Errorf("unknown acknowlegement policy %v", p)
	}
}

type ReplayPolicy int

const (
	ReplayInstant ReplayPolicy = iota
	ReplayOriginal
)

func (p ReplayPolicy) String() string {
	switch p {
	case ReplayInstant:
		return "Instant"
	case ReplayOriginal:
		return "Original"
	default:
		return "Unknown Replay Policy"
	}
}

func (p *ReplayPolicy) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case jsonString("instant"):
		*p = ReplayInstant
	case jsonString("original"):
		*p = ReplayOriginal
	default:
		return fmt.Errorf("can not unmarshal %q", data)
	}

	return nil
}

func (p ReplayPolicy) MarshalJSON() ([]byte, error) {
	switch p {
	case ReplayOriginal:
		return json.Marshal("original")
	case ReplayInstant:
		return json.Marshal("instant")
	default:
		return nil, fmt.Errorf("unknown replay policy %v", p)
	}
}

var (
	AckAck      = []byte("+ACK")
	AckNak      = []byte("-NAK")
	AckProgress = []byte("+WPI")
	AckNext     = []byte("+NXT")
	AckTerm     = []byte("+TERM")
)

type DeliverPolicy int

const (
	DeliverAll DeliverPolicy = iota
	DeliverLast
	DeliverNew
	DeliverByStartSequence
	DeliverByStartTime
	DeliverLastPerSubject
)

func (p DeliverPolicy) String() string {
	switch p {
	case DeliverAll:
		return "All"
	case DeliverLast:
		return "Last"
	case DeliverNew:
		return "New"
	case DeliverByStartSequence:
		return "By Start Sequence"
	case DeliverByStartTime:
		return "By Start Time"
	case DeliverLastPerSubject:
		return "Last Per Subject"
	default:
		return "Unknown Deliver Policy"
	}
}

func (p *DeliverPolicy) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case jsonString("all"), jsonString("undefined"):
		*p = DeliverAll
	case jsonString("last"):
		*p = DeliverLast
	case jsonString("new"):
		*p = DeliverNew
	case jsonString("by_start_sequence"):
		*p = DeliverByStartSequence
	case jsonString("by_start_time"):
		*p = DeliverByStartTime
	case jsonString("last_per_subject"):
		*p = DeliverLastPerSubject
	}

	return nil
}

func (p DeliverPolicy) MarshalJSON() ([]byte, error) {
	switch p {
	case DeliverAll:
		return json.Marshal("all")
	case DeliverLast:
		return json.Marshal("last")
	case DeliverNew:
		return json.Marshal("new")
	case DeliverByStartSequence:
		return json.Marshal("by_start_sequence")
	case DeliverByStartTime:
		return json.Marshal("by_start_time")
	case DeliverLastPerSubject:
		return json.Marshal("last_per_subject")
	default:
		return nil, fmt.Errorf("unknown deliver policy %v", p)
	}
}

// PriorityPolicy determines policy for selecting messages based on priority.
type PriorityPolicy int

const (
	PriorityNone PriorityPolicy = iota
	PriorityOverflow
	PriorityPinnedClient
)

func (p PriorityPolicy) String() string {
	switch p {
	case PriorityOverflow:
		return "Overflow"
	case PriorityPinnedClient:
		return "Pinned Client"
	default:
		return "None"
	}
}

func (p *PriorityPolicy) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case jsonString("none"):
		*p = PriorityNone
	case jsonString("overflow"):
		*p = PriorityOverflow
	case jsonString("pinned_client"):
		*p = PriorityPinnedClient
	default:
		return fmt.Errorf("unknown priority policy: %v", string(data))
	}
	return nil
}

func (p PriorityPolicy) MarshalJSON() ([]byte, error) {
	switch p {
	case PriorityNone:
		return json.Marshal("none")
	case PriorityOverflow:
		return json.Marshal("overflow")
	case PriorityPinnedClient:
		return json.Marshal("pinned_client")
	default:
		return nil, fmt.Errorf("unknown priority policy: %v", p)
	}
}

// ConsumerConfig is the configuration for a JetStream consumes
//
// NATS Schema Type io.nats.jetstream.api.v1.consumer_configuration
type ConsumerConfig struct {
	Description        string          `json:"description,omitempty"`
	AckPolicy          AckPolicy       `json:"ack_policy"`
	AckWait            time.Duration   `json:"ack_wait,omitempty"`
	DeliverPolicy      DeliverPolicy   `json:"deliver_policy"`
	DeliverSubject     string          `json:"deliver_subject,omitempty"`
	DeliverGroup       string          `json:"deliver_group,omitempty"`
	Durable            string          `json:"durable_name,omitempty"` // Durable is deprecated. All consumers will have names. picked by clients.
	Name               string          `json:"name,omitempty"`
	FilterSubject      string          `json:"filter_subject,omitempty"`
	FilterSubjects     []string        `json:"filter_subjects,omitempty"`
	FlowControl        bool            `json:"flow_control,omitempty"`
	Heartbeat          time.Duration   `json:"idle_heartbeat,omitempty"`
	MaxAckPending      int             `json:"max_ack_pending,omitempty"`
	MaxDeliver         int             `json:"max_deliver,omitempty"`
	BackOff            []time.Duration `json:"backoff,omitempty"`
	MaxWaiting         int             `json:"max_waiting,omitempty"`
	OptStartSeq        uint64          `json:"opt_start_seq,omitempty"`
	OptStartTime       *time.Time      `json:"opt_start_time,omitempty"`
	RateLimit          uint64          `json:"rate_limit_bps,omitempty"`
	ReplayPolicy       ReplayPolicy    `json:"replay_policy"`
	SampleFrequency    string          `json:"sample_freq,omitempty"`
	HeadersOnly        bool            `json:"headers_only,omitempty"`
	MaxRequestBatch    int             `json:"max_batch,omitempty"`
	MaxRequestExpires  time.Duration   `json:"max_expires,omitempty"`
	MaxRequestMaxBytes int             `json:"max_bytes,omitempty"`
	InactiveThreshold  time.Duration   `json:"inactive_threshold,omitempty"`
	Replicas           int             `json:"num_replicas"`
	MemoryStorage      bool            `json:"mem_storage,omitempty"`
	// Metadata is additional metadata for the Consumer.
	Metadata map[string]string `json:"metadata,omitempty"`

	// PauseUntil is for suspending the consumer until the deadline.
	PauseUntil time.Time `json:"pause_until,omitempty"`

	// Priority groups
	PriorityGroups []string       `json:"priority_groups,omitempty"`
	PriorityPolicy PriorityPolicy `json:"priority_policy,omitempty"`
	PinnedTTL      time.Duration  `json:"priority_timeout,omitempty"`

	// Don't add to general clients.
	Direct bool `json:"direct,omitempty"`
}

// SequenceInfo is the consumer and stream sequence that uniquely identify a message
type SequenceInfo struct {
	Consumer uint64     `json:"consumer_seq"`
	Stream   uint64     `json:"stream_seq"`
	Last     *time.Time `json:"last_active,omitempty"`
}

// ConsumerInfo reports the current state of a consumer
type ConsumerInfo struct {
	Stream         string               `json:"stream_name"`
	Name           string               `json:"name"`
	Config         ConsumerConfig       `json:"config"`
	Created        time.Time            `json:"created"`
	Delivered      SequenceInfo         `json:"delivered"`
	AckFloor       SequenceInfo         `json:"ack_floor"`
	NumAckPending  int                  `json:"num_ack_pending"`
	NumRedelivered int                  `json:"num_redelivered"`
	NumWaiting     int                  `json:"num_waiting"`
	NumPending     uint64               `json:"num_pending"`
	Cluster        *ClusterInfo         `json:"cluster,omitempty"`
	PushBound      bool                 `json:"push_bound,omitempty"`
	Paused         bool                 `json:"paused,omitempty"`
	PauseRemaining time.Duration        `json:"pause_remaining,omitempty"`
	TimeStamp      time.Time            `json:"ts"`
	PriorityGroups []PriorityGroupState `json:"priority_groups,omitempty"`
}

// PriorityGroupState is the state of a consumer group
type PriorityGroupState struct {
	Group          string    `json:"group"`
	PinnedClientID string    `json:"pinned_client_id,omitempty"`
	PinnedTS       time.Time `json:"pinned_ts,omitempty"`
}

// JSApiConsumerGetNextRequest is for getting next messages for pull based consumers
//
// NATS Schema Type io.nats.jetstream.api.v1.consumer_getnext_request
type JSApiConsumerGetNextRequest struct {
	Expires   time.Duration `json:"expires,omitempty"`
	Batch     int           `json:"batch,omitempty"`
	MaxBytes  int           `json:"max_bytes,omitempty"`
	NoWait    bool          `json:"no_wait,omitempty"`
	Heartbeat time.Duration `json:"idle_heartbeat,omitempty"`
}

// ConsumerNakOptions is for optional NAK parameters, e.g. delay.
type ConsumerNakOptions struct {
	Delay time.Duration `json:"delay"`
}

func jsonString(s string) string {
	return "\"" + s + "\""
}
