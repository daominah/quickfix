package quickfix

import "time"

// initMsgSeqNum is MsgSeqNum after store_Reset.
// HNX fucker wants MsgSeqNum start at 0, the whole world want MsgSeqNum start at 1
const initMsgSeqNum = 0

//The MessageStore interface provides methods to record and retrieve messages for resend purposes
type MessageStore interface {
	// NextSenderMsgSeqNum only reads
	NextSenderMsgSeqNum() int
	// NextTargetMsgSeqNum only reads
	NextTargetMsgSeqNum() int

	IncrNextSenderMsgSeqNum() error
	IncrNextTargetMsgSeqNum() error

	SetNextSenderMsgSeqNum(next int) error
	SetNextTargetMsgSeqNum(next int) error

	CreationTime() time.Time

	SaveMessage(seqNum int, msg []byte) error
	GetMessages(beginSeqNum, endSeqNum int) ([][]byte, error)

	Refresh() error
	Reset() error

	Close() error
}

//The MessageStoreFactory interface is used by session to create a session specific message store
type MessageStoreFactory interface {
	Create(sessionID SessionID) (MessageStore, error)
}

type memoryStore struct {
	senderMsgSeqNum, targetMsgSeqNum int
	creationTime                     time.Time
	// Difference HNX messages can have same seqNum CON CAC
	messageMap                       map[int][][]byte
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		senderMsgSeqNum: initMsgSeqNum,
		targetMsgSeqNum: initMsgSeqNum,
	}
}

func (store *memoryStore) NextSenderMsgSeqNum() int {
	return store.senderMsgSeqNum + 1
}

func (store *memoryStore) NextTargetMsgSeqNum() int {
	return store.targetMsgSeqNum + 1
}

func (store *memoryStore) IncrNextSenderMsgSeqNum() error {
	store.senderMsgSeqNum++
	return nil
}

func (store *memoryStore) IncrNextTargetMsgSeqNum() error {
	store.targetMsgSeqNum++
	return nil
}

func (store *memoryStore) SetNextSenderMsgSeqNum(nextSeqNum int) error {
	store.senderMsgSeqNum = nextSeqNum - 1
	return nil
}
func (store *memoryStore) SetNextTargetMsgSeqNum(nextSeqNum int) error {
	store.targetMsgSeqNum = nextSeqNum - 1
	return nil
}

func (store *memoryStore) CreationTime() time.Time {
	return store.creationTime
}

func (store *memoryStore) Reset() error {
	store.senderMsgSeqNum = initMsgSeqNum
	store.targetMsgSeqNum = initMsgSeqNum
	store.creationTime = time.Now()
	store.messageMap = nil
	return nil
}

func (store *memoryStore) Refresh() error {
	//nop, nothing to refresh
	return nil
}

func (store *memoryStore) Close() error {
	//nop, nothing to close
	return nil
}

func (store *memoryStore) SaveMessage(seqNum int, msg []byte) error {
	if store.messageMap == nil {
		store.messageMap = make(map[int][][]byte)
	}

	store.messageMap[seqNum] =  append(store.messageMap[seqNum], msg)
	return nil
}

func (store *memoryStore) GetMessages(beginSeqNum, endSeqNum int) ([][]byte, error) {
	var msgs [][]byte
	for seqNum := beginSeqNum; seqNum <= endSeqNum; seqNum++ {
		if ms, ok := store.messageMap[seqNum]; ok {
			msgs = append(msgs, ms...)
		}
	}
	return msgs, nil
}

type memoryStoreFactory struct{}

func (f memoryStoreFactory) Create(sessionID SessionID) (MessageStore, error) {
	m := newMemoryStore()
	m.Reset()
	return m, nil
}

//NewMemoryStoreFactory returns a MessageStoreFactory instance that created in-memory MessageStores
func NewMemoryStoreFactory() MessageStoreFactory { return memoryStoreFactory{} }
