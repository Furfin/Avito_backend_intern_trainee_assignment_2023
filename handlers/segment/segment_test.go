package segment

import (
	"bytes"
	"example/ravito/initializers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostSegmentError(t *testing.T) {
	body := []byte(`{
		"slug": "seg1"
	}`)
	initializers.LoadEnvVars()
	initializers.SetupLogger()
	initializers.ConnectToDB()
	req, err := http.NewRequest("POST", "/segment", bytes.NewReader(body))

	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	handler := http.HandlerFunc(CreateSegment)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestPostSegmentCreate(t *testing.T) {
	body := []byte(`{
		"slug": "segNew",
		"upadd":1
	}`)
	initializers.LoadEnvVars()
	initializers.SetupLogger()
	initializers.ConnectToDB()
	req, err := http.NewRequest("POST", "/segment", bytes.NewReader(body))

	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	handler := http.HandlerFunc(CreateSegment)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestPostSegmentDelete(t *testing.T) {
	body := []byte(`{
		"slug": "segNew"
	}`)
	initializers.LoadEnvVars()
	initializers.SetupLogger()
	initializers.ConnectToDB()
	req, err := http.NewRequest("DELETE", "/segment", bytes.NewReader(body))

	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	handler := http.HandlerFunc(DeleteSegment)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
