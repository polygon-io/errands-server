
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
		TTL 				int 		`json:"ttl,omitempty"`
		Retries 			int 		`json:"retries,omitempty"`
		Priority 			int 		`json:"priority,omitempty"`
		DeleteOnCompleted 	bool 		`json:"deleteOnCompleted,omitempty"`
	} 							`json:"options"`
	Data 			*gin.H 		`json:"data,omitempty"`
	Created 		int64 		`json:"created"`
	Status 			string 		`json:"status,omitempty"`

	// Internal attributes:
	Progress		float64 	`json:"progress"`
	Attempts 		int 		`json:"attempts"`
	Started			int64 		`json:"started,omitempty"` // Timestamp of last Start
	Failed			int64 		`json:"failed,omitempty"` // Timestamp of last Fail
	Completed		int64 		`json:"compelted,omitempty"` // Timestamp of last Fail
	Logs 			[]Log 		`json:"logs"`
}


var LogSeverities []string = []string{ "INFO", "WARNING", "ERROR" }
type Log struct {
	Severity 		string 		`json:"severity" binding:"required"`
	Message 		string 		`json:"message" binding:"required"`
	Timestamp 		int64 		`json:"timestamp"`
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
	e.Logs = make( []Log, 0 )
}


func ( e *Errand ) addToLogs( severity, message string ) error {
	obj := Log{
		Severity: severity,
		Message: message,
		Timestamp: getTimestamp(),
	}
	e.Logs = append( e.Logs, obj )
	return nil
}

