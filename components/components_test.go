package components

import (
	"testing"
)

func Test_ConfParse(t *testing.T) {
	err := ParseConfig("../dev.conf.toml")
	if err != nil {
		t.Errorf("parse config error:%v", err.Error())
		return
	}
	t.Log(config)
}
