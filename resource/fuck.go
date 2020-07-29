package resource

import (
	"github.com/luchanghe/jx3Robot/tool"
)

func SearchFuck() string {
	msg := tool.HttpGet(fuckYou)
	return msg
}
