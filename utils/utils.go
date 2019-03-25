
package utils


import (
	"time"
)


func GetTimestamp() int64 {
	return ( time.Now().UnixNano() / 1000000 )
}


func Contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
