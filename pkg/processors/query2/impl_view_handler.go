/*
 * Copyright (c) 2025-present unTill Software Development Group B.V.
 * @author Michael Saigachenko
 */
package query2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/voedger/voedger/pkg/appdef"
	"github.com/voedger/voedger/pkg/coreutils"
	"github.com/voedger/voedger/pkg/istructs"
	"github.com/voedger/voedger/pkg/istructsmem"
	"github.com/voedger/voedger/pkg/pipeline"
)

type viewHandler struct {
}

var _ IApiPathHandler = (*viewHandler)(nil) // ensure that viewHandler implements IApiPathHandler

func (h *viewHandler) CheckRateLimit(ctx context.Context, qw *queryWork) error {
	// TODO: implement rate limits check
	return nil
}

func (h *viewHandler) CheckType(ctx context.Context, qw *queryWork) error {
	switch qw.iWorkspace {
	case nil:
		// workspace is dummy
		if qw.iView = appdef.View(qw.appStructs.AppDef().Type, qw.msg.QName()); qw.iView == nil {
			return coreutils.NewHTTPErrorf(http.StatusBadRequest, fmt.Sprintf("view %s does not exist", qw.msg.QName()))
		}
	default:
		if qw.iView = appdef.View(qw.iWorkspace.Type, qw.msg.QName()); qw.iView == nil {
			return coreutils.NewHTTPErrorf(http.StatusBadRequest, fmt.Sprintf("view %s does not exist in %v", qw.msg.QName(), qw.iWorkspace))
		}
	}
	return nil
}

func (h *viewHandler) ResultType(ctx context.Context, qw *queryWork, statelessResources istructsmem.IStatelessResources) error {
	qw.resultType = qw.iView
	return nil
}

func (h *viewHandler) AuthorizeRequest(ctx context.Context, qw *queryWork) error {
	// ws := qw.iWorkspace
	// if ws == nil {
	// 	// workspace is dummy
	// 	ws = qw.iView.Workspace()
	// }
	// ok, err := qw.appPart.IsOperationAllowed(ws, appdef.OperationKind_Select, qw.msg.QName(), nil, qw.roles)
	// if err != nil {
	// 	return err
	// }
	// if !ok {
	// 	return coreutils.NewHTTPError(http.StatusForbidden, errors.New(""))
	// }
	return nil
}

func (h *viewHandler) AuthorizeResult(ctx context.Context, qw *queryWork) error {
	// if qw.iQuery.Result() != appdef.AnyType {
	// 	// will authorize result only if result is sys.Any
	// 	// otherwise each field is considered as allowed if EXECUTE ON QUERY is allowed
	// 	return nil
	// }
	// ws := qw.iWorkspace
	// if ws == nil {
	// 	// workspace is dummy
	// 	ws = qw.iQuery.Workspace()
	// }
	// for _, elem := range qw.queryParams.Elements() {
	// 	nestedPath := elem.Path().AsArray()
	// 	nestedType := qw.resultType
	// 	for _, nestedName := range nestedPath {
	// 		if len(nestedName) == 0 {
	// 			// root
	// 			continue
	// 		}
	// 		// incorrectness is excluded already on validation stage in [queryParams.validate]
	// 		containersOfNested := nestedType.(appdef.IWithContainers)
	// 		// container presence is checked already on validation stage in [queryParams.validate]
	// 		nestedContainer := containersOfNested.Container(nestedName)
	// 		nestedType = nestedContainer.Type()
	// 	}
	// 	requestedfields := []string{}
	// 	for _, resultField := range elem.ResultFields() {
	// 		requestedfields = append(requestedfields, resultField.Field())
	// 	}

	// 	// TODO: temporary solution. To be eliminated after implementing ACL in VSQL for Air
	// 	ok := oldacl.IsOperationAllowed(appdef.OperationKind_Select, nestedType.QName(), requestedfields, oldacl.EnrichPrincipals(qw.principals, qw.msg.WSID()))
	// 	if !ok {
	// 		if ok, err = qw.appPart.IsOperationAllowed(ws, appdef.OperationKind_Select, nestedType.QName(), requestedfields, qw.roles); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	if !ok {
	// 		return coreutils.NewSysError(http.StatusForbidden)
	// 	}
	// }
	return nil
}
func (h *viewHandler) RowsProcessor(ctx context.Context, qw *queryWork) (err error) {
	err = h.validateFields(qw)
	if err != nil {
		return
	}
	oo := make([]*pipeline.WiredOperator, 0)
	if len(qw.queryParams.Constraints.Order) != 0 || qw.queryParams.Constraints.Skip > 0 || qw.queryParams.Constraints.Limit > 0 {
		oo = append(oo, pipeline.WireAsyncOperator("Aggregator", newAggregator(qw.queryParams)))
	}
	o, err := newFilter(qw)
	if err != nil {
		return
	}
	oo = append(oo, pipeline.WireAsyncOperator("Filter", o))
	if len(qw.queryParams.Constraints.Keys) != 0 {
		oo = append(oo, pipeline.WireAsyncOperator("Keys", newKeys(qw.queryParams.Constraints.Keys)))
	}
	oo = append(oo, pipeline.WireAsyncOperator("Sender", &sender{responder: qw.msg.Responder()}))
	qw.rowsProcessor = pipeline.NewAsyncPipeline(ctx, "View rows processor", oo[0], oo[1:]...)
	return
}
func (h *viewHandler) Exec(ctx context.Context, qw *queryWork) (err error) {
	kk, err := h.getKeys(qw)
	if err != nil {
		return
	}
	for i := range kk {
		err = qw.appStructs.ViewRecords().Read(ctx, qw.msg.WSID(), kk[i], func(key istructs.IKey, value istructs.IValue) (err error) {
			obj := objectBackedByMap{}
			obj.data = coreutils.FieldsToMap(key, qw.appStructs.AppDef())
			for k, v := range coreutils.FieldsToMap(value, qw.appStructs.AppDef()) {
				obj.data[k] = v
			}
			return qw.callbackFunc(obj)
		})
		if err != nil {
			return
		}
	}
	return
}
func (h *viewHandler) getKeys(qw *queryWork) (keys []istructs.IKeyBuilder, err error) {
	for _, field := range qw.appStructs.AppDef().Type(qw.iView.QName()).(appdef.IView).Key().Fields() {
		intf := qw.queryParams.Constraints.Where[field.Name()]
		switch intf.(type) {
		case map[string]interface{}:
			if len(keys) != 0 {
				return
			}
			in, ok := intf.(map[string]interface{})["$in"]
			if !ok {
				return nil, errUnsupportedConstraint
			}
			params, ok := in.([]interface{})
			if !ok {
				return nil, errUnexpectedParams
			}
			for i := range params {
				switch params[i].(type) {
				case json.Number:
					key := qw.appStructs.ViewRecords().KeyBuilder(qw.iView.QName())
					key.PutNumber(field.Name(), params[i].(json.Number))
					keys = append(keys, key)
				case string:
					key := qw.appStructs.ViewRecords().KeyBuilder(qw.iView.QName())
					key.PutString(field.Name(), params[i].(string))
					keys = append(keys, key)
				default:
					return nil, errUnsupportedType
				}
			}
		case json.Number:
			if len(keys) == 0 {
				key := qw.appStructs.ViewRecords().KeyBuilder(qw.iView.QName())
				key.PutNumber(field.Name(), intf.(json.Number))
				keys = append(keys, key)
			} else {
				keys[0].PutNumber(field.Name(), intf.(json.Number))
			}
		case string:
			if len(keys) == 0 {
				key := qw.appStructs.ViewRecords().KeyBuilder(qw.iView.QName())
				key.PutString(field.Name(), intf.(string))
				keys = append(keys, key)
			} else {
				keys[0].PutString(field.Name(), intf.(string))
			}
		case nil:
			return
		default:
			return nil, errUnsupportedType
		}
	}

	return
}
func (h *viewHandler) validateFields(qw *queryWork) (err error) {
	view := qw.appStructs.AppDef().Type(qw.iView.QName()).(appdef.IView)

	if qw.queryParams.Constraints == nil {
		return errConstraintsAreNull
	}
	if len(qw.queryParams.Constraints.Where) == 0 {
		return errWhereConstraintIsEmpty
	}
	if _, ok := qw.queryParams.Constraints.Where[view.Key().Fields()[0].Name()]; !ok {
		return errWhereConstraintMustHaveFirstMemberOfThePartitionKey
	}

	ff := make(map[string]bool)
	for _, field := range view.Fields() {
		ff[field.Name()] = true
	}
	for k := range qw.queryParams.Constraints.Where {
		if !ff[k] {
			return fmt.Errorf("%w: '%s'", errUnexpectedField, k)
		}
	}
	return
}
