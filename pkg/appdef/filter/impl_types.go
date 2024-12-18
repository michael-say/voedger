/*
 * Copyright (c) 2024-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package filter

import (
	"fmt"
	"iter"

	"github.com/voedger/voedger/pkg/appdef"
	"github.com/voedger/voedger/pkg/goutils/set"
)

// typesFilter is a filter that matches types by their kind.
//
// # Supports:
//   - appdef.IFilter.
//   - fmt.Stringer
type typesFilter struct {
	filter
	ws    appdef.QName
	types appdef.TypeKindSet
}

func makeTypesFilter(ws appdef.QName, t appdef.TypeKind, tt ...appdef.TypeKind) appdef.IFilter {
	f := &typesFilter{ws: ws, types: set.From(t)}
	f.types.Set(tt...)
	return f
}

func (typesFilter) Kind() appdef.FilterKind { return appdef.FilterKind_Types }

func (f typesFilter) Match(t appdef.IType) bool {
	if !f.types.Contains(t.Kind()) {
		return false
	}

	if f.ws == appdef.NullQName {
		return true
	}

	ws := t.Workspace()
	return (ws != nil) && (ws.QName() == f.ws)
}

func (f typesFilter) String() string {
	var s string
	if t, ok := typesStringDecorators[string(f.types.AsBytes())]; ok {
		s = t
	} else {
		// TYPES(…)
		// TYPES(…) FROM …)
		s = "TYPES("
		for i, t := range f.types.All() {
			if i > 0 {
				s += ", "
			}
			s += t.TrimString()
		}
		s += ")"
	}
	if f.ws != appdef.NullQName {
		s += fmt.Sprintf(" FROM %s", f.ws)
	}
	return s
}

func (f typesFilter) Types() iter.Seq[appdef.TypeKind] { return f.types.Values() }

var typesStringDecorators = map[string]string{
	string(appdef.TypeKind_Structures.AsBytes()): "ALL TABLES",
	string(appdef.TypeKind_Functions.AsBytes()):  "ALL FUNCTIONS",
}
