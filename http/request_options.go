package http

// RequestOptions is the optional arguments of one API call, as a value.
type RequestOptions struct {
	// PathParams are values for the {name} placeholders in the path template.
	PathParams map[string]any
	// QueryParams are normalised by the value table before they are sent.
	QueryParams map[string]any
	// JSONBody is a JSON-encodable value, raw JSON text sent as is, or a DTOPayload; not with FormData or Files.
	JSONBody any
	// FormData is form-encoded data to send.
	FormData map[string]any
	// Files are multipart files to send.
	Files map[string]any
	// Headers are extra request headers, merged case-insensitively; an authenticated call replaces Authorization.
	Headers map[string]string
	// Unauthenticated omits the bearer token when true (default false = authenticated).
	Unauthenticated bool
	// RaiseOnError overrides SdkConfig.RaiseOnError per call; nil defers to the config.
	RaiseOnError *bool
	// ResponseFormat overrides SdkConfig.ResponseFormat for this call; "" defers to the config.
	ResponseFormat string
	// DTOShape is the typed view this call produces for DTO(); a raw verb call leaves it nil.
	DTOShape DTOShape
}

// Frozen copies every reference-typed field, so the call's content is fixed before caller code runs.
func (o RequestOptions) Frozen() RequestOptions {
	out := o
	out.PathParams = snapshotAnyMap(o.PathParams)
	out.QueryParams = snapshotAnyMap(o.QueryParams)
	out.FormData = snapshotAnyMap(o.FormData)
	out.Files = snapshotAnyMap(o.Files)
	out.Headers = snapshotStringMap(o.Headers)
	out.JSONBody = snapshotValue(o.JSONBody)
	return out
}

// ReportedBody is the body this call carries, for printing and logging; files are summarised apart.
func (o RequestOptions) ReportedBody() any {
	if o.JSONBody != nil {
		return o.JSONBody
	}
	if o.FormData != nil {
		return o.FormData
	}
	return nil
}
