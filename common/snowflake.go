package common

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

var (
	Snowflake *snowflake.Node
)

func init() {
	var err error
	Snowflake, err = snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}
}
