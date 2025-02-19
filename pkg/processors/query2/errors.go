/*
 * Copyright (c) 2025-present unTill Pro, Ltd.
 *
 * @author Daniil Solovyov
 */

package query2

import "errors"

var (
	errConstraintsAreNull                             = errors.New("constraints are null")
	errWhereConstraintIsEmpty                         = errors.New("where constraint is empty")
	errWhereConstraintMustHaveFirstMemberOfPrimaryKey = errors.New("where constraint must have first member of key")
	errUnsupportedConstraint                          = errors.New("unsupported constraint")
	errUnexpectedParams                               = errors.New("unexpected params")
	errUnsupportedType                                = errors.New("unsupported type")
	errUnexpectedField                                = errors.New("unexpected field")
)
