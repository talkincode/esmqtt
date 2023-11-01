package zaplog

import (
	"testing"

	"github.com/talkincode/esmqtt/common/zaplog/log"
)

func TestInfo(t *testing.T) {
	log.Info("test")
}
