
package main


import (
	gin "github.com/gin-gonic/gin"
)






/*

Name
Type
TTL
Retries
Priority
Data

 */


type Errand struct {

	// Marshalling Attributes:
	ID 				string 		`json:"id"`
	Name 			string 		`json:"name" binding:"required"`
	Type 			string 		`json:"type" binding:"required"`
	Options 		struct {
		TTL 			int 		`json:"ttl,omitempty"`
		Retries 		int 		`json:"retries,omitempty"`
		Priority 		int 		`json:"priority,omitempty"`
	} 							`json:"options"`
	Data 			*gin.H 		`json:"data,omitempty"`
	Created 		int64 		`json:"created"`
	Status 			string 		`json:"status,omitempty"`

	// Internal attributes:
	

}


