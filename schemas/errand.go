// Schemas provides model schemas for the errand server.
package schemas

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/polygon-io/errands-server/utils"
)

type Status string

var (
	StatusBlocked   Status = "blocked"
	StatusInactive  Status = "inactive"
	StatusActive    Status = "active"
	StatusFailed    Status = "failed"
	StatusCompleted Status = "completed"
)

var ErrandStatuses = []Status{StatusBlocked, StatusInactive, StatusActive, StatusFailed, StatusCompleted}

//easyjson:json
type Errand struct {

	// General Attributes:
	ID      string `json:"id"`
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Options struct {
		TTL               int  `json:"ttl,omitempty"`
		Retries           int  `json:"retries,omitempty"`
		Priority          int  `json:"priority,omitempty"`
		DeleteOnCompleted bool `json:"deleteOnCompleted,omitempty"`
	} `json:"options"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Created int64                  `json:"created"`
	Status  Status                 `json:"status,omitempty"`
	Results map[string]interface{} `json:"results,omitempty"`

	// Internal attributes:
	Progress  float64 `json:"progress"`
	Attempts  int     `json:"attempts"`
	Started   int64   `json:"started,omitempty"`   // Timestamp of last Start
	Failed    int64   `json:"failed,omitempty"`    // Timestamp of last Fail
	Completed int64   `json:"compelted,omitempty"` // Timestamp of last Fail
	Logs      []Log   `json:"logs,omitempty"`

	// PipelineID is the ID of the pipeline that this errand belongs to (if any)
	PipelineID string `json:"pipeline,omitempty"`
}

var LogSeverities = []string{"INFO", "WARNING", "ERROR"}

//easyjson:json
type Log struct {
	Severity  string `json:"severity" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Timestamp int64  `json:"timestamp"`
}

func NewErrand() *Errand {
	obj := &Errand{}
	obj.SetDefaults()

	return obj
}

func (e *Errand) SetDefaults() {
	uid := uuid.New()

	uidText, err := uid.MarshalText()
	if err != nil {
		panic(err)
	}

	e.ID = string(uidText)
	e.Status = "inactive"
	e.Created = utils.GetTimestamp()
	e.Logs = make([]Log, 0)
}

func (e *Errand) AddToLogs(severity, message string) error {
	if !utils.Contains(LogSeverities, severity) {
		return fmt.Errorf("invalid log severity: %s", severity)
	}

	obj := Log{
		Severity:  severity,
		Message:   message,
		Timestamp: utils.GetTimestamp(),
	}
	e.Logs = append(e.Logs, obj)

	return nil
}
