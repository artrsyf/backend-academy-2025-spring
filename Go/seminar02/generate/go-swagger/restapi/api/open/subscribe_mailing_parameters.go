// Code generated by go-swagger; DO NOT EDIT.

package open

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewSubscribeMailingParams creates a new SubscribeMailingParams object
//
// There are no default values defined in the spec.
func NewSubscribeMailingParams() SubscribeMailingParams {

	return SubscribeMailingParams{}
}

// SubscribeMailingParams contains all the bound params for the subscribe mailing operation
// typically these are obtained from a http.Request
//
// swagger:parameters subscribeMailing
type SubscribeMailingParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: query
	*/
	Email strfmt.Email
	/*
	  Required: true
	  In: query
	*/
	Name string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSubscribeMailingParams() beforehand.
func (o *SubscribeMailingParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qEmail, qhkEmail, _ := qs.GetOK("email")
	if err := o.bindEmail(qEmail, qhkEmail, route.Formats); err != nil {
		res = append(res, err)
	}

	qName, qhkName, _ := qs.GetOK("name")
	if err := o.bindName(qName, qhkName, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindEmail binds and validates parameter Email from query.
func (o *SubscribeMailingParams) bindEmail(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("email", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("email", "query", raw); err != nil {
		return err
	}

	// Format: email
	value, err := formats.Parse("email", raw)
	if err != nil {
		return errors.InvalidType("email", "query", "strfmt.Email", raw)
	}
	o.Email = *(value.(*strfmt.Email))

	if err := o.validateEmail(formats); err != nil {
		return err
	}

	return nil
}

// validateEmail carries on validations for parameter Email
func (o *SubscribeMailingParams) validateEmail(formats strfmt.Registry) error {

	if err := validate.FormatOf("email", "query", "email", o.Email.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindName binds and validates parameter Name from query.
func (o *SubscribeMailingParams) bindName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("name", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("name", "query", raw); err != nil {
		return err
	}
	o.Name = raw

	return nil
}
