package api_helpers

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strconv"
)

func RenderErr(w http.ResponseWriter, errs ...*jsonapi.ErrorObject) {
	if len(errs) == 0 {
		panic("expected non-empty errors slice")
	}

	// getting status of first occurred error
	status, err := strconv.ParseInt(errs[0].Status, 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse status"))
	}
	w.Header().Set("content-type", jsonapi.MediaType)
	w.WriteHeader(int(status))
	jsonapi.MarshalErrors(w, errs)
}

func Render(w http.ResponseWriter, res interface{}) {
	w.Header().Set("content-type", jsonapi.MediaType)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(errors.Wrap(err, "failed to render response"))
	}
}

type BadRequester interface {
	BadRequest() map[string]error
}

func BadRequest(err error) []*jsonapi.ErrorObject {
	cause := errors.Cause(err)
	if cause == io.EOF {
		return []*jsonapi.ErrorObject{
			{
				Title:  http.StatusText(http.StatusBadRequest),
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Detail: "Request body were expected",
			},
		}
	}

	switch cause := cause.(type) {
	case validation.Errors:
		return toJsonapiErrors(cause)
	case BadRequester:
		return toJsonapiErrors(cause.BadRequest())
	default:
		return []*jsonapi.ErrorObject{
			{
				Title:  http.StatusText(http.StatusBadRequest),
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Detail: "Your request was invalid in some way",
			},
		}
	}

}

func toJsonapiErrors(m map[string]error) []*jsonapi.ErrorObject {
	errs := make([]*jsonapi.ErrorObject, 0, len(m))
	for key, value := range m {
		errs = append(errs, &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusBadRequest),
			Status: fmt.Sprintf("%d", http.StatusBadRequest),
			Meta: &map[string]interface{}{
				"field": key,
				"error": value.Error(),
			},
		})
	}
	return errs
}
