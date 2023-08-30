package user

import (
	"bytes"
	"context"
	"example/ravito/initializers"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestPostUserError(t *testing.T) {
	initializers.LoadEnvVars()
	initializers.SetupLogger()
	initializers.ConnectToDB()
	req := httptest.NewRequest(http.MethodPost, "/user/:userid/add", bytes.NewReader(nil))
	ctx := req.Context()
	ctx = context.WithValue(ctx, httprouter.ParamsKey, httprouter.Params{
		{"userid", "-1"},
	})
	req = req.WithContext(ctx)
	fmt.Println(httprouter.ParamsFromContext(req.Context()).ByName("userid"))

	handler := http.HandlerFunc(UserSegmentsUpdate)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
