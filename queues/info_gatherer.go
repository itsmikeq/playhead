package queues

import (
	"bytes"
	"encoding/json"
	"sync"
	"wtc_go/wtc"
)

// used to gather info and store it in memory
// then create a buffer to be written to s3

var buffer bytes.Buffer
var mu = &sync.Mutex{}

type AllInfo struct {
	Dpps []wtc.DynamoPlayedPositionS `json:"play_positions"`
	Lpps []wtc.DynamoLastPlayed      `json:"last_playeds"`
}

func withLockContext(fn func()) {
	mu.Lock()
	defer mu.Unlock()

	fn()
}

// function is to take in bytes and return a json blob
func GetAllItems(userUUID string) (AllInfo, error) {
	var err error
	var allInfo AllInfo
	// maps are not thread safe, so lock it up
	withLockContext(func() { allInfo.Dpps, err = wtc.GetPlayedPositionsToString(userUUID) })
	withLockContext(func() { allInfo.Lpps, err = wtc.GetLastPlayeds(userUUID) })
	if ErrorHandler(err) {
		return allInfo, err
	}
	return allInfo, nil
}

func JsonifyAllInfo(allInfo AllInfo) string {
	if js, err := json.Marshal(allInfo); ErrorHandler(err) {
		return ""
	} else {
		return string(js)
	}
}
