package uuid

import (
	"github.com/donetkit/gin-contrib/utils/snowflake"
	"log"
)

var snowFlakeNode *snowflake.Node

func init() {
	var err error
	snowFlakeNode, err = snowflake.NewNode(0)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func SnowFlakeId() int64 {
	return snowFlakeNode.Generate().Int64()
}
