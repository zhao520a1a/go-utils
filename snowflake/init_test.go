package snowflake

import (
	"fmt"
	"testing"
)

func TestMsgID(t *testing.T) {
	fmt.Printf("msgID %d", MsgID())
}
