/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 */

package appdef_test

import (
	"fmt"
	"strings"

	"github.com/voedger/voedger/pkg/appdef"
)

func ExampleIAppDefBuilder_AddData() {

	var app appdef.IAppDef

	numName := appdef.NewQName("test", "num")
	floatName := appdef.NewQName("test", "float")
	strName := appdef.NewQName("test", "string")
	tokenName := appdef.NewQName("test", "token")
	weekDayName := appdef.NewQName("test", "weekDay")

	// how to build AppDef with data types
	{
		appDef := appdef.New()

		appDef.AddData(numName, appdef.DataKind_int64, appdef.NullQName, appdef.MinExcl(0)).SetComment("Natural number")

		_ = appDef.AddData(floatName, appdef.DataKind_float64, appdef.NullQName)

		_ = appDef.AddData(strName, appdef.DataKind_string, appdef.NullQName, appdef.MinLen(1), appdef.MaxLen(4))

		_ = appDef.AddData(tokenName, appdef.DataKind_string, strName, appdef.Pattern("^[a-z]+$"))

		_ = appDef.AddData(weekDayName, appdef.DataKind_string, strName, appdef.Enum("Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"))

		if a, err := appDef.Build(); err == nil {
			app = a
		} else {
			panic(err)
		}
	}

	// how to inspect data types in builded AppDef
	{
		cnt := 0
		app.DataTypes(false, func(d appdef.IData) {
			cnt++
			fmt.Println("-", d, "inherits from", d.Ancestor())
			if d.Comment() != "" {
				fmt.Println(" ", d.Comment())
			}
			str := []string{}
			d.Constraints(func(c appdef.IConstraint) {
				str = append(str, fmt.Sprint(c))
			})
			if len(str) > 0 {
				fmt.Printf("  constraints: (%v)\n", strings.Join(str, `, `))
			}
		})
		fmt.Println("overall data types: ", cnt)
	}

	// Output:
	// - float64-data «test.float» inherits from float64-data «sys.float64»
	// - int64-data «test.num» inherits from int64-data «sys.int64»
	//   Natural number
	//   constraints: (MinExcl: 0)
	// - string-data «test.string» inherits from string-data «sys.string»
	//   constraints: (MinLen: 1, MaxLen: 4)
	// - string-data «test.token» inherits from string-data «test.string»
	//   constraints: (Pattern: `^[a-z]+$`)
	// - string-data «test.weekDay» inherits from string-data «test.string»
	//   constraints: (Enum: [Fri Mon Sat Sun Thu Tue Wed])
	// overall data types:  5
}