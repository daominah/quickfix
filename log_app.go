package quickfix

import (
	"github.com/daominah/gomicrokit/log"
	"github.com/daominah/quickfixdev"
)

// ZLogger implements the Log interface by using the `gomicrokit/log`
type ZLogger struct {
	// sessionID of global log is the empty string
	sessionID string
}

func (l ZLogger) OnIncoming(msg []byte) {
	log.Debugf("FIXLog incoming %v: %s", l.sessionID, msg)
}
func (l ZLogger) OnOutgoing(msg []byte) {
	log.Debugf("FIXLog outgoing %v: %s", l.sessionID, msg)
}
func (l ZLogger) OnEvent(event string) {
	log.Debugf("FIXLog event %v: %v", l.sessionID, event)
}
func (l ZLogger) OnEventf(eventTemplate string, args ...interface{}) {
	args = append([]interface{}{l.sessionID}, args...)
	log.Debugf("FIXLog event %v: "+eventTemplate, args...)
}

// ZLogFactory implements the LogFactory by using the `gomicrokit/log`
type ZLogFactory struct{}

func NewZLogFactory() *ZLogFactory { return &ZLogFactory{} }
func (f ZLogFactory) Create() (Log, error) {
	return ZLogger{sessionID: ""}, nil
}
func (f ZLogFactory) CreateSessionLog(sessionID SessionID) (Log, error) {
	return ZLogger{sessionID: sessionID.String()}, nil
}

// LogApp implements Application interface.
// It logs human readable FIX messages by the `gomicrokit/log`
type LogApp struct{}

func NewLogApp() *LogApp { return &LogApp{} }
func (a LogApp) OnCreate(sessionID SessionID) {
	log.Infof("created a session %v", sessionID.String())
}
func (a LogApp) OnLogon(sessionID SessionID) {
	log.Infof("session %v logged on", sessionID.String())
}
func (a LogApp) OnLogout(sessionID SessionID) {
	log.Infof("session %v logged out or disconnected", sessionID.String())
}
func (a LogApp) FromAdmin(msg *Message, sessionID SessionID) MessageRejectError {
	//if msg.IsMsgTypeOf("0") { // Heartbeat
	//	return nil
	//}
	log.Infof("received FromAdmin: %v", quickfixdev.HumanString(msg.String()))
	return nil
}
func (a LogApp) FromApp(
	msg *Message, sessionID SessionID) MessageRejectError {
	log.Infof("received FromApp: %v", quickfixdev.HumanString(msg.String()))
	return nil
}
func (a LogApp) ToAdmin(msg *Message, sessionID SessionID) {
	//if msg.IsMsgTypeOf("0") { // Heartbeat
	//	return
	//}
	log.Infof("send ToAdmin: %v", quickfixdev.HumanString(msg.String()))
}
func (a LogApp) ToApp(msg *Message, std SessionID) error {
	log.Infof("send ToApp: %v", quickfixdev.HumanString(msg.String()))
	return nil
}
