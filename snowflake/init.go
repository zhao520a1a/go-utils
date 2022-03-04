package snowflake

import (
	"math/rand"
	"time"
)

var (
	node *Node
)

func init() {
	// 此处先用随机值
	rand.Seed(time.Now().UnixNano())
	var err error
	node, err = NewNode(int64(rand.Uint64() % 1023))
	if err != nil {
		panic(err)
	}
}

// MsgID 获得消息id
func MsgID() int64 {
	return node.Generate().Int64()
}
