package render

import (
	"github.com/benpate/data"
	"github.com/benpate/data/option"
	"github.com/benpate/datatype"
	"github.com/benpate/exp"
)

// ModelService interface wraps the generic Object-* functions that standard services provide
type ModelService interface {
	ObjectNew() data.Object
	ObjectList(exp.Expression, ...option.Option) (data.Iterator, error)
	ObjectLoad(exp.Expression) (data.Object, error)
	ObjectSave(data.Object, string) error
	ObjectDelete(data.Object, string) error
	Debug() datatype.Map
}
