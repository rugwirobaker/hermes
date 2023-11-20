package hermes

import "context"

type Environment interface {
	// Type returns the name of the environment type.
	Type() string

	// GetNodeRole returns the role of the node.
	GetNodeRole(context.Context) (string, error)
}
