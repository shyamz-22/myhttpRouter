package assert

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func ResponseWithBody(t *testing.T, w *httptest.ResponseRecorder, expectedStatusCode int, expectedResponseBody string) {
	ResponseWithStatus(t, w, expectedStatusCode)

	body, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {

		t.Fatalf("unexpected error occured: %v\n", err)
	}
	if expectedResponseBody != string(body) {
		t.Fatalf("\nExpected: %s\nActual:%s\n", expectedResponseBody, string(body))
	}
}

func ResponseWithStatus(t *testing.T, w *httptest.ResponseRecorder, expectedStatusCode int) {
	if w.Result().StatusCode != expectedStatusCode {
		t.Fatalf("\nExpected: %d\nActual:%d\n", expectedStatusCode, w.Result().StatusCode)
	}
}
