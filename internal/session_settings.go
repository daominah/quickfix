package internal

import "time"

//SessionSettings stores all of the configuration for a given session
type SessionSettings struct {
	ResetOnLogon                 bool
	RefreshOnLogon               bool
	ResetOnLogout                bool
	ResetOnDisconnect            bool
	HeartBtInt                   time.Duration
	SessionTime                  *TimeRange
	// InitiateLogon is true means we are client
	InitiateLogon                bool
	ResendRequestChunkSize       int
	EnableLastMsgSeqNumProcessed bool
	SkipCheckLatency             bool
	MaxLatency                   time.Duration
	DisableMessagePersist        bool

	//required on logon for FIX.T.1 messages
	DefaultApplVerID string

	//specific to initiators
	ReconnectInterval    time.Duration
	LogoutTimeout        time.Duration
	LogonTimeout         time.Duration
	SocketConnectAddress []string

	// HNXVersion is specific to HaNoiStockExchange
	HNXVersion string
}
