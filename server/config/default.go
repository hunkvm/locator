package config

import (
	"time"

	"github.com/hunkvm/locator/pkg/security"
)

func defaultServerConfig() ServerConfig {
	return ServerConfig{
		Listeners: ListenersConfig{
			Client: "127.0.0.1:8881",
			Peer:   "127.0.0.1:9991",
		},
		ClientChannelSecurity: security.ServerTLSConfig{
			TLSPolicy: security.TLSPolicyNone,
		},
		PeerChannelSecurity: security.PeerTLSConfig{
			TLSPolicy: security.TLSPolicyNone,
		},
		// TODO: Must check value of raft with default value
		Raft: RaftConfig{
			NodeID:         1,
			PeerAddrs:      "1@127.0.0.1:9991",
			ElectionTick:   10,
			Heartbeat:      100 * time.Millisecond,
			SnapshotCount:  100,
			CatchupEntries: 0,
		},
		// TODO: Must check value of health check with default value
		HealthCheck: HealthCheckConfig{
			MinBackoffDuration: 500 * time.Millisecond,
			MaxBackoffDuration: 5 * time.Minute,
			MinConnectTimeout:  100 * time.Millisecond,
			MaxBackoffAttempts: 5,
		},
		DataDir: "/var/lib/locator",
		Logger: LogConfig{
			Level:  LogLevelInfo,
			Format: LogFormatText,
		},
	}
}
