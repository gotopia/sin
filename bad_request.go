package sin

import (
	"fmt"
	"strings"

	edpb "google.golang.org/genproto/googleapis/rpc/errdetails"
)

// BadRequest describes violations in a client request.
type BadRequest struct {
	fvs []*FieldViolation
}

// NewBadRequest returns a bad request error.
func NewBadRequest(fvs ...*FieldViolation) *BadRequest {
	return &BadRequest{
		fvs: fvs,
	}
}

// WithFieldViolations returns a new BadRequest with the provided field violations appended to the bad request.
func (r *BadRequest) WithFieldViolations(fvs ...*FieldViolation) *BadRequest {
	r.fvs = append(r.fvs, fvs...)
	return r
}

// Serialize BadRequest to proto message.
func (r *BadRequest) Serialize() *edpb.BadRequest {
	var fvpbs []*edpb.BadRequest_FieldViolation
	for _, fv := range r.fvs {
		fvpbs = append(fvpbs, fv.Serialize())
	}
	return &edpb.BadRequest{
		FieldViolations: fvpbs,
	}
}

// FieldViolation describes a single bad request field.
type FieldViolation struct {
	field       string
	description string
}

// NewFieldViolationWithRule returns a field violation with rule.
func NewFieldViolationWithRule(field string, rule string, params ...string) *FieldViolation {
	formatter, ok := ruleToFormatter[rule]
	if !ok {
		formatter = ruleToFormatter["invalid"]
	}
	return &FieldViolation{
		field:       field,
		description: formatter(field, params...),
	}
}

// NewFieldViolation returns a field violation.
func NewFieldViolation(field string, description string) *FieldViolation {
	return &FieldViolation{
		field:       field,
		description: description,
	}
}

// Serialize FieldViolation to proto message.
func (v *FieldViolation) Serialize() *edpb.BadRequest_FieldViolation {
	return &edpb.BadRequest_FieldViolation{
		Field:       v.field,
		Description: v.description,
	}
}

type formatter func(field string, params ...string) string

var ruleToFormatter = map[string]formatter{
	"required": func(field string, params ...string) string {
		return fmt.Sprintf("%v can't be blank", field)
	},
	"max": func(field string, params ...string) string {
		count := params[0]
		if count == "1" {
			return fmt.Sprintf("%v is too long (maximum is %v character)", field, count)
		}
		return fmt.Sprintf("%v is too long (maximum is %v character)", field, count)
	},
	"min": func(field string, params ...string) string {
		count := params[0]
		if count == "1" {
			return fmt.Sprintf("%v is too short (minimum is %v characters", field, count)
		}
		return fmt.Sprintf("%v is too short (minimum is %v characters)", field, count)
	},
	"numeric": func(field string, params ...string) string {
		return fmt.Sprintf("%v is not a number", field)
	},
	"gt": func(field string, params ...string) string {
		threshold := params[0]
		return fmt.Sprintf("%v must be greater than %v", field, threshold)
	},
	"gte": func(field string, params ...string) string {
		threshold := params[0]
		return fmt.Sprintf("%v must be greater than or equal to %v", field, threshold)
	},
	"eq": func(field string, params ...string) string {
		threshold := params[0]
		return fmt.Sprintf("%v must be equal to %v", field, threshold)
	},
	"lt": func(field string, params ...string) string {
		threshold := params[0]
		return fmt.Sprintf("%v must be less than %v", field, threshold)
	},
	"lte": func(field string, params ...string) string {
		threshold := params[0]
		return fmt.Sprintf("%v must be less than or equal to %v", field, threshold)
	},
	"oneof": func(field string, params ...string) string {
		arrayString := params[0]
		array := strings.Split(arrayString, " ")
		return fmt.Sprintf("%v is not included in the list %v", field, array)
	},
	"invalid": func(field string, params ...string) string {
		return fmt.Sprintf("%v is invalid", field)
	},
}
