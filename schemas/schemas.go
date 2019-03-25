
package schemas

import (
	"errors"
	// gin "github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
	utils "github.com/polygon-io/errands-server/utils"
)





var ErrandStatuses []string = []string{"inactive", "active", "failed", "completed"}
//easyjson:json
type Errand struct {

	// General Attributes:
	ID 				string 			`json:"id"`
	Name 			string 			`json:"name" binding:"required"`
	Type 			string 			`json:"type" binding:"required"`
	Options 		struct {
		TTL 				int 		`json:"ttl,omitempty"`
		Retries 			int 		`json:"retries,omitempty"`
		Priority 			int 		`json:"priority,omitempty"`
		DeleteOnCompleted 	bool 		`json:"deleteOnCompleted,omitempty"`
	} 											`json:"options"`
	Data 			map[string]interface{} 		`json:"data,omitempty"`
	Created 		int64 			`json:"created"`
	Status 			string 			`json:"status,omitempty"`
	Results 		map[string]interface{} 		`json:"results,omitempty"`

	// Internal attributes:
	Progress		float64 		`json:"progress"`
	Attempts 		int 			`json:"attempts"`
	Started			int64 			`json:"started,omitempty"`		// Timestamp of last Start
	Failed			int64 			`json:"failed,omitempty"` 		// Timestamp of last Fail
	Completed		int64 			`json:"compelted,omitempty"` 	// Timestamp of last Fail
	Logs 			[]Log 			`json:"logs,omitempty"`
}


var LogSeverities []string = []string{ "INFO", "WARNING", "ERROR" }
//easyjson:json
type Log struct {
	Severity 		string 			`json:"severity" binding:"required"`
	Message 		string 			`json:"message" binding:"required"`
	Timestamp 		int64 			`json:"timestamp"`
}




func NewErrand() *Errand {
	obj := &Errand{}
	obj.SetDefaults()
	return obj
}



func ( e *Errand ) SetDefaults(){
	uid := uuid.New()
	uidText, err := uid.MarshalText(); if err != nil {
		panic( err )
	}
	e.ID = string( uidText )
	e.Status = "inactive"
	e.Created = utils.GetTimestamp()
	e.Logs = make( []Log, 0 )
}


func ( e *Errand ) AddToLogs( severity, message string ) error {
	if !utils.Contains( LogSeverities, severity ) {
		return errors.New("Invalid log severity")
	}
	obj := Log{
		Severity: severity,
		Message: message,
		Timestamp: utils.GetTimestamp(),
	}
	e.Logs = append( e.Logs, obj )
	return nil
}


