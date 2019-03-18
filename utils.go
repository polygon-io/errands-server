
package main

import (
	"time"
)


func getTimestamp() int64 {
	return ( time.Now().UnixNano() / 1000000 )
}





