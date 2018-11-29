package sin

import (
	"fmt"

	"github.com/pkg/errors"
	edpb "google.golang.org/genproto/googleapis/rpc/errdetails"
)

// DebugInfo describes additional debugging info.
type DebugInfo struct {
	st     errors.StackTrace
	detail string
}

// NewDebugInfo returns a debug info.
func NewDebugInfo(err error) *DebugInfo {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	var st errors.StackTrace
	if stackTracer, ok := err.(stackTracer); ok {
		st = stackTracer.StackTrace()
	}
	return &DebugInfo{
		st:     st,
		detail: err.Error(),
	}
}

// Serialize DebugInfo to proto message.
func (i *DebugInfo) Serialize() *edpb.DebugInfo {
	var stackEntries []string
	for _, e := range i.st {
		stackEntries = append(stackEntries, fmt.Sprintf("%+v", e))
	}
	return &edpb.DebugInfo{
		StackEntries: stackEntries,
		Detail:       i.detail,
	}
}
