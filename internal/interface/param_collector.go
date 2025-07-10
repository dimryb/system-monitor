package interfaces

import "context"

type ParamCollector interface {
	Collect(context.Context) (string, error)
}
