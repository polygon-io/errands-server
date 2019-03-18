
package main


import (
	uuid "github.com/google/uuid"
	gin "github.com/gin-gonic/gin"
)





type Errand struct {

	// General Attributes:
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
	Progress		int16 		`json:"progress"`
	Attempts 		int16 		`json:"attempts"`
	Started			int64 		`json:"started,omitempty"`
	FailedReason	string 		`json:"failedReason,omitempty"`

}



func NewErrand() *Errand {
	obj := &Errand{}
	obj.setDefaults()
	return obj
}



func ( e *Errand ) setDefaults(){
	uid := uuid.New()
	uidText, err := uid.MarshalText(); if err != nil {
		panic( err )
	}
	e.ID = string( uidText )
	e.Status = "inactive"
	e.Created = getTimestamp()
}

