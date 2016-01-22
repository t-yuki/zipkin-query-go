package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/strfmt"
)

// NewGetTraceTraceIDParams creates a new GetTraceTraceIDParams object
// with the default values initialized.
func NewGetTraceTraceIDParams() GetTraceTraceIDParams {
	var ()
	return GetTraceTraceIDParams{}
}

// GetTraceTraceIDParams contains all the bound params for the get trace trace ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetTraceTraceID
type GetTraceTraceIDParams struct {
	/*the 64-bit hex-encoded id of the trace as a path parameter.
	  Required: true
	  In: path
	*/
	TraceID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetTraceTraceIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	rTraceID, rhkTraceID, _ := route.Params.GetOK("traceId")
	if err := o.bindTraceID(rTraceID, rhkTraceID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTraceTraceIDParams) bindTraceID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.TraceID = raw

	return nil
}