package husky

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/rand"
)

type UUID interface {
	Id() Int64
	Uuid() string
	TimeID() string
}

type _UUID struct {
	v *snowflake.Node
}

func (u *_UUID) Id() Int64 {
	return Int64(u.v.Generate().Int64())
}

func (u *_UUID) Uuid() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

func (u *_UUID) TimeID() string {
	now := time.Now().Format("20060102150405")
	n := rand.Int63n(1000000000000000000)
	return fmt.Sprintf("%s%018d", now, n)
}

var _UuidIns map[string]UUID

func init() {
	_UuidIns = make(map[string]UUID)
}

func InitUuid(nodeSeq int64, key ...string) {
	if nodeSeq > 1023 {
		panic("nodeSeq must < 1024")
	}
	node, err := snowflake.NewNode(nodeSeq)
	if err != nil {
		panic(err)
	}
	_ins := &_UUID{node}
	if len(key) == 0 {
		_UuidIns[""] = _ins
	} else {
		_UuidIns[key[0]] = _ins
	}
}

func Uuid(key ...string) UUID {
	if len(key) == 0 {
		return _UuidIns[""]
	} else {
		return _UuidIns[key[0]]
	}
}
