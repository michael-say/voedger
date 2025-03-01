/*
 * Copyright (c) 2024-present unTill Software Development Group B.V.
 * @author Michael Saigachenko
 */

package sys_it

import (
	"bytes"
	"testing"

	"github.com/voedger/voedger/pkg/appdef"
	"github.com/voedger/voedger/pkg/goutils/testingu/require"
	"github.com/voedger/voedger/pkg/istructs"
	it "github.com/voedger/voedger/pkg/vit"
)

type apiWithTypes struct {
	appdef.IAppDef
}

func (a *apiWithTypes) PublishedTypes(ws appdef.IWorkspace, role appdef.QName, callback appdef.PublishedTypeCallback) error {
	callback(a.Type(appdef.NewQName("app1pkg", "articles")), []appdef.OperationKind{appdef.OperationKind_Insert, appdef.OperationKind_Update})
	callback(a.Type(appdef.NewQName("app1pkg", "department")), []appdef.OperationKind{appdef.OperationKind_Insert, appdef.OperationKind_Update})
	return nil
}

func TestOpenApi(t *testing.T) {
	vit := it.NewVIT(t, &it.SharedConfig_App1)
	defer vit.TearDown()
	require := require.New(t)

	//ws := vit.WS(istructs.AppQName_test1_app1, "test_ws")

	appDef, err := vit.AppDef(istructs.AppQName_test1_app1)
	ws := appDef.Workspace(appdef.NewQName("app1pkg", "test_wsWS"))
	require.NotNil(ws)
	require.NoError(err)

	apiDef2 := &apiWithTypes{appDef}
	writer := new(bytes.Buffer)
	err = appdef.GenOpenAPISchema(apiDef2, ws, appdef.NewQName("app1pkg", "admin"), writer, "Test API", "1.0")

	require.NoError(err)
	t.Log(writer.String())

	require.Equal("", writer.String())

}
