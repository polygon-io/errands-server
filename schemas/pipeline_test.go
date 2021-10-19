package schemas

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipelineValidate(t *testing.T) {
	t.Run("no errands", func(t *testing.T) {
		p := Pipeline{Name: "no errands!"}
		assert.Error(t, p.Validate())
	})

	t.Run("errands with duplicate name", func(t *testing.T) {
		p := Pipeline{
			Name: "duplicate name",
			Errands: []*Errand{{
				Name: "errand",
			}, {
				Name: "errand",
			}},
		}

		assert.Error(t, p.Validate())
	})

	t.Run("errand with self dependency", func(t *testing.T) {
		p := Pipeline{
			Name:    "self dependency",
			Errands: []*Errand{{Name: "single errand"}},
			Dependencies: []*PipelineDependency{{
				Target:    "single errand",
				DependsOn: "single errand",
			}},
		}

		assert.Error(t, p.Validate())
	})

	t.Run("dependency with invalid target", func(t *testing.T) {
		p := Pipeline{
			Name:    "self dependency",
			Errands: []*Errand{{Name: "single errand"}},
			Dependencies: []*PipelineDependency{{
				Target:    "not a real errand",
				DependsOn: "single errand",
			}},
		}

		assert.Error(t, p.Validate())
	})

	t.Run("dependency with invalid dependsOn", func(t *testing.T) {
		p := Pipeline{
			Name:    "self dependency",
			Errands: []*Errand{{Name: "single errand"}},
			Dependencies: []*PipelineDependency{{
				Target:    "single errand",
				DependsOn: "not a real dependency",
			}},
		}

		assert.Error(t, p.Validate())
	})

	t.Run("simple dependency cycle | 2 errands", func(t *testing.T) {
		/*
			A <--> B
		*/
		p := Pipeline{
			Name: "simple dependency cycle",
			Errands: []*Errand{
				{Name: "A"},
				{Name: "B"},
			},
			Dependencies: []*PipelineDependency{
				{Target: "A", DependsOn: "B"},
				{Target: "B", DependsOn: "A"},
			},
		}

		assert.Error(t, p.Validate())
	})

	t.Run("strongly connected subgraph cycle", func(t *testing.T) {
		/*
			A <--> B    C
		*/
		p := Pipeline{
			Name: "strongly connected cycle",
			Errands: []*Errand{
				{Name: "A"},
				{Name: "B"},
				{Name: "C"},
			},
			Dependencies: []*PipelineDependency{
				{Target: "A", DependsOn: "B"},
				{Target: "B", DependsOn: "A"},
			},
		}

		assert.Error(t, p.Validate())
	})

	t.Run("single graph with cycle", func(t *testing.T) {
		/*
			A <--> B <-- C
		*/
		p := Pipeline{
			Name: "single graph with cycle",
			Errands: []*Errand{
				{Name: "A"},
				{Name: "B"},
				{Name: "C"},
			},
			Dependencies: []*PipelineDependency{
				{Target: "A", DependsOn: "B"},
				{Target: "B", DependsOn: "C"},
				{Target: "B", DependsOn: "A"},
			},
		}

		assert.Error(t, p.Validate())
	})

	t.Run("multiple sub-graphs one with cycle", func(t *testing.T) {
		/*
			A <--> B <-- C    D --> E --> F
		*/
		p := Pipeline{
			Name: "strongly connected cycle",
			Errands: []*Errand{
				{Name: "A"},
				{Name: "B"},
				{Name: "C"},
				{Name: "D"},
				{Name: "E"},
				{Name: "F"},
			},
			Dependencies: []*PipelineDependency{
				{Target: "A", DependsOn: "B"},
				{Target: "B", DependsOn: "C"},
				{Target: "B", DependsOn: "A"},

				{Target: "E", DependsOn: "D"},
				{Target: "F", DependsOn: "E"},
			},
		}

		assert.Error(t, p.Validate())
	})

	t.Run("multiple sub-graphs happy path", func(t *testing.T) {
		/*
			A --> B --> C    D --> E --> F
		*/
		p := Pipeline{
			Name: "strongly connected cycle",
			Errands: []*Errand{
				{Name: "A"},
				{Name: "B"},
				{Name: "C"},
				{Name: "D"},
				{Name: "E"},
				{Name: "F"},
			},
			Dependencies: []*PipelineDependency{
				{Target: "A", DependsOn: "B"},
				{Target: "B", DependsOn: "C"},

				{Target: "E", DependsOn: "D"},
				{Target: "F", DependsOn: "E"},
			},
		}

		assert.NoError(t, p.Validate())
	})

	t.Run("single graph happy path | diverging", func(t *testing.T) {
		/*
			    |--> D
			A --|
			    |--> B --> C
		*/
		p := Pipeline{
			Name: "single graph with cycle",
			Errands: []*Errand{
				{Name: "A"},
				{Name: "B"},
				{Name: "C"},
				{Name: "D"},
			},
			Dependencies: []*PipelineDependency{
				{Target: "B", DependsOn: "A"},
				{Target: "C", DependsOn: "B"},
				{Target: "D", DependsOn: "A"},
			},
		}

		assert.NoError(t, p.Validate())
	})

	t.Run("single graph happy path | converging", func(t *testing.T) {
		/*
			A --> B --|
			          |--> C --> E
			      D --|
		*/
		p := Pipeline{
			Name: "single graph with cycle",
			Errands: []*Errand{
				{Name: "A"},
				{Name: "B"},
				{Name: "C"},
				{Name: "D"},
				{Name: "E"},
			},
			Dependencies: []*PipelineDependency{
				{Target: "B", DependsOn: "A"},
				{Target: "C", DependsOn: "B"},
				{Target: "C", DependsOn: "D"},
				{Target: "E", DependsOn: "C"},
			},
		}

		assert.NoError(t, p.Validate())
	})

	t.Run("single graph happy path | converging and diverging", func(t *testing.T) {
		/*
		        |--> B --|        |--> E
		   A -->|        |--> C --|
		        |--> D --|        |--> F --> G
		*/
		p := Pipeline{
			Name: "single graph with cycle",
			Errands: []*Errand{
				{Name: "A"},
				{Name: "B"},
				{Name: "C"},
				{Name: "D"},
				{Name: "E"},
				{Name: "F"},
				{Name: "G"},
			},
			Dependencies: []*PipelineDependency{
				{Target: "D", DependsOn: "A"},
				{Target: "B", DependsOn: "A"},
				{Target: "C", DependsOn: "B"},
				{Target: "C", DependsOn: "D"},
				{Target: "E", DependsOn: "C"},
				{Target: "F", DependsOn: "C"},
				{Target: "G", DependsOn: "F"},
			},
		}

		assert.NoError(t, p.Validate())
	})
}
