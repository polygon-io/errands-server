package pipelines

import "github.com/polygon-io/errands-server/schemas"

type PipelineDependency struct {
	// Target is the name of the errand within the pipeline that this dependency relates to
	Target string `json:"target" binding:"required"`

	// DependsOn is the name of the errand within the pipeline that the Target errand depends on
	DependsOn string `json:"depends_on" binding:"required"`
}

type Pipeline struct {
	Name              string `json:"name" binding:"required"`
	DeleteOnCompleted bool   `json:"deleteOnCompleted,omitempty"`

	Errands      []schemas.Errand     `json:"errands" binding:"required"`
	Dependencies []PipelineDependency `json:"dependencies,omitempty"`
}
