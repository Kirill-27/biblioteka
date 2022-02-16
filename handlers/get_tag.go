package handlers

import (
	"fmt"
	"github.com/Kirill-27/biblioteka/api_helpers"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func GetTag(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		api_helpers.RenderErr(w, api_helpers.BadRequest(err)...)
	}
	fmt.Println(id)
}
