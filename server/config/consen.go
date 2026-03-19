package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/hunkvm/locator/pkg/validation"
)

type RaftConfig struct {
	NodeID         uint64
	PeerAddrs      string
	ElectionTick   int
	Heartbeat      time.Duration
	SnapshotCount  uint64
	CatchupEntries uint64
}

func (c RaftConfig) Validate() []validation.Error {
	errors := make([]validation.Error, 0, 6)

	if c.NodeID == 0 {
		errors = append(errors, validation.Error{
			Field: "NodeID", Value: c.NodeID,
			Reason: "must be positive integer",
		})
	}

	errors = append(errors, validatePeerAddrs(c.PeerAddrs)...)

	if err := validation.ValidateInteger(c.ElectionTick,
		new(0), new(1000)); err != nil {
		errors = append(errors, validation.Error{
			Field: "ElectionTick", Value: c.ElectionTick,
			Reason: err.Error(),
		})
	}

	if err := validation.ValidateDuration(c.Heartbeat,
		new(0*time.Millisecond), new(30*time.Minute)); err != nil {
		errors = append(errors, validation.Error{
			Field: "Heartbeat", Value: c.Heartbeat,
			Reason: err.Error(),
		})
	}

	if err := validation.ValidateInteger(c.SnapshotCount,
		new(uint64(0)), new(uint64(1_000_000))); err != nil {
		errors = append(errors, validation.Error{
			Field: "SnapshotCount", Value: c.SnapshotCount,
			Reason: err.Error(),
		})
	}

	if err := validation.ValidateInteger(c.CatchupEntries,
		new(uint64(0)), new(uint64(1_000_000))); err != nil {
		errors = append(errors, validation.Error{
			Field: "CatchupEntries", Value: c.CatchupEntries,
			Reason: err.Error(),
		})
	}

	return errors
}

func validatePeerAddrs(peer string) []validation.Error {
	var errors []validation.Error
	parts := strings.SplitN(peer, "@", 2)
	if len(parts) != 2 {
		errors = append(errors, validation.Error{
			Field: "PeerAddrs", Value: peer,
			Reason: "expected format <id>@<host>:<port>",
		})
		return errors
	}
	id, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil || id == 0 {
		errors = append(errors, validation.Error{
			Field: "PeerAddrs", Value: peer,
			Reason: "node id must be a positive integer",
		})
	}
	if parts[1] == "" {
		errors = append(errors, validation.Error{
			Field: "PeerAddrs", Value: peer,
			Reason: "peer address can not be empty",
		})
	}
	return errors
}
