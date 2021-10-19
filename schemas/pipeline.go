package schemas

import (
	"errors"
	"fmt"

	"github.com/polygon-io/errands-server/utils"
)

// Pipeline represents an errand pipeline which consists of errands and dependencies between them.
// Cyclic dependencies are not allowed.
//easyjson:json
type Pipeline struct {
	Name              string `json:"name" binding:"required"`
	DeleteOnCompleted bool   `json:"deleteOnCompleted,omitempty"`

	Errands      []*Errand             `json:"errands,omitempty" binding:"required"`
	Dependencies []*PipelineDependency `json:"dependencies,omitempty"`

	// Attributes added by errands server
	ID     string `json:"id"`
	Status Status `json:"status,omitempty"`

	// TODO: use ptime once that's public
	StartedMillis int64 `json:"startedMillis"`
	EndedMillis   int64 `json:"endedMillis,omitempty"`
}

// PipelineDependency describes a dependency between errands within a pipeline.
//easyjson:json
type PipelineDependency struct {
	// Target is the name of the errand within the pipeline that this dependency relates to
	Target string `json:"target" binding:"required"`

	// DependsOn is the name of the errand within the pipeline that the Target errand depends on
	DependsOn string `json:"dependsOn" binding:"required"`
}

type dependencyGraph struct {
	// errands maps errand names to the errands themselves for quick look-ups
	errands map[string]*Errand

	// dependencyToDependents maps an errand to a slice of errands that depend on it.
	// ie dependencyToDependents["A"] = []{"B"} -> "B" depends on "A"
	dependencyToDependents map[string][]*Errand

	// dependentToDependencies maps an errand to a slice of errands that it depends on.
	// This is the transpose of dependencyToDependents.
	// ie dependentToDependencies["B"] = []{"A"} -> "B" depends on "A"
	dependentToDependencies map[string][]*Errand
}

// Validate checks that the pipeline describes a valid dependency graph and returns a user-friendly error if the pipeline is invalid.
func (p *Pipeline) Validate() error {
	if len(p.Errands) == 0 {
		return errors.New("no errands specified in pipeline")
	}

	// Map of errand name to errand schema
	errandMap := make(map[string]*Errand, len(p.Errands))

	// Build up map of errand name to errand, and double check there are no errands with duplicate names
	for _, errand := range p.Errands {
		if _, exists := errandMap[errand.Name]; exists {
			return fmt.Errorf("duplicate name in errands list: %s", errand.Name)
		}

		errandMap[errand.Name] = errand
	}

	// Make sure dependencies all reference valid errands in the pipeline
	for _, dep := range p.Dependencies {
		if dep.Target == dep.DependsOn {
			return fmt.Errorf("errand cannot depend on itself: %s", dep.Target)
		}

		if _, exists := errandMap[dep.Target]; !exists {
			return fmt.Errorf("dependency references unknown errand name: %s", dep.Target)
		}

		if _, exists := errandMap[dep.DependsOn]; !exists {
			return fmt.Errorf("dependency references unknown errand name: %s", dep.DependsOn)
		}
	}

	// Check for dependency cycles
	return p.buildDependencyGraph().checkForDependencyCycles()
}

func (p *Pipeline) GetUnblockedErrands() []*Errand {
	return p.buildDependencyGraph().findUnblockedErrands()
}

func (p *Pipeline) RecalculateStatus() {
	var numCompleted int

	for _, errand := range p.Errands {
		// Failed takes precedence over everything
		if errand.Status == StatusFailed {
			p.Status = StatusFailed
			break
		}

		if p.Status == StatusInactive && errand.Status == StatusActive {
			p.Status = StatusActive
		}

		if errand.Status == StatusCompleted {
			numCompleted++
		}
	}

	// If all the errands are completed, the pipeline is complete
	if numCompleted == len(p.Errands) {
		p.Status = StatusCompleted

		// Update the ended timestamp if it wasn't set already
		if p.EndedMillis == 0 {
			p.EndedMillis = utils.GetTimestamp()
		}
	}
}

// buildDependencyGraph constructs a dependencyGraph for this pipeline.
// This function assumes the errands and dependencies have been validated already.
func (p *Pipeline) buildDependencyGraph() dependencyGraph {
	g := dependencyGraph{
		errands:                 make(map[string]*Errand, len(p.Errands)),
		dependencyToDependents:  make(map[string][]*Errand, len(p.Errands)),
		dependentToDependencies: make(map[string][]*Errand, len(p.Errands)),
	}

	// Build up the errand map and initialize the graph with all errands and no dependencies
	for _, errand := range p.Errands {
		g.errands[errand.Name] = errand
		g.dependencyToDependents[errand.Name] = []*Errand{}
		g.dependentToDependencies[errand.Name] = []*Errand{}
	}

	// Fill in the dependencies
	for _, dep := range p.Dependencies {
		g.dependencyToDependents[dep.DependsOn] = append(g.dependencyToDependents[dep.DependsOn], g.errands[dep.Target])
		g.dependentToDependencies[dep.Target] = append(g.dependentToDependencies[dep.Target], g.errands[dep.DependsOn])
	}

	return g
}

// findUnblockedErrands returns a slice of errands that have no dependencies blocking them.
func (g dependencyGraph) findUnblockedErrands() []*Errand {
	var independentErrands []*Errand

	for dependent, dependencies := range g.dependentToDependencies {
		var blocked bool

		// This dependent is blocked if any of its dependencies are not completed
		for _, dependency := range dependencies {
			if dependency.Status != StatusCompleted {
				blocked = true
				break
			}
		}

		if !blocked {
			independentErrands = append(independentErrands, g.errands[dependent])
		}
	}

	return independentErrands
}

// checkForDependencyCycles returns an error if it finds a cyclic dependency in the dependencyGraph.
// It determines cycles by doing the following:
// For each node N, perform a depth first traversal of G starting at N.
// If we ever encounter N again during the traversal, we've detected a cycle.
func (g dependencyGraph) checkForDependencyCycles() error {
	independentErrands := g.findUnblockedErrands()
	if len(independentErrands) == 0 {
		return fmt.Errorf("dependency cycle found; all errands have dependencies")
	}

	allErrands := make([]*Errand, 0, len(g.errands))
	for _, errand := range g.errands {
		allErrands = append(allErrands, errand)
	}

	// Prime the visit stack with the current home errand (N).
	currentHomeErrand := allErrands[0]
	toVisitStack := []*Errand{currentHomeErrand}
	nextIndex := 1

	// currentTreeVisitedSet keeps track of the errands we've seen in the tree traversal starting at 'currentHomeErrand'
	// This set gets cleared every time we finish a traversal and move on to the next node.
	currentTreeVisitedSet := make(map[string]struct{}, len(g.dependencyToDependents))

	for len(toVisitStack) > 0 {
		topOfStackIndex := len(toVisitStack) - 1
		errand := toVisitStack[topOfStackIndex]
		toVisitStack = toVisitStack[:topOfStackIndex] // Pop off the last value from the stack

		// If we've made it back to the current path's start index, we've found a cycle.
		if _, exists := currentTreeVisitedSet[errand.Name]; exists && errand.Name == currentHomeErrand.Name {
			return fmt.Errorf("dependency cycle found involving '%s'", errand.Name)
		}

		// Add this errand to the visited set
		currentTreeVisitedSet[errand.Name] = struct{}{}

		// Add all of this errand's dependencies to the visit stack
		toVisitStack = append(toVisitStack, g.dependencyToDependents[errand.Name]...)

		// If our visit stack is empty, we've exhausted the nodes in this tree without finding a cycle.
		// If there are more nodes to traverse from, reset and go again.
		if len(toVisitStack) == 0 && nextIndex < len(allErrands) {
			currentHomeErrand = allErrands[nextIndex]
			toVisitStack = []*Errand{currentHomeErrand}
			currentTreeVisitedSet = make(map[string]struct{})
			nextIndex++
		}
	}

	return nil
}
