package ch

import (
	"fmt"
)

var (
	// ErrGraphIsFrozen Graph is frozen, so it can not be modified.
	ErrGraphIsFrozen                    = fmt.Errorf("Graph has been frozen")
	ErrSourceAndTargetListsCanNotBeSame = fmt.Errorf("Source and target lists cannot be same")
	ErrVertexNotFoundInGraph            = fmt.Errorf("Vertex not found in graph")
)
