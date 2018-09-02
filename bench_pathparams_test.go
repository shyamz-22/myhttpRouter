package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkWithPathParamOneRequest(b *testing.B) {

	rtr := New()
	rtr.Add("/repos/:owner/:repo/pulls/:number/merge", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/repos/shyamz-22/oidc/pulls/44/merge", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rtr.ServeHTTP(w, r)
	}
}

func BenchmarkWithPathParamTenRequests(b *testing.B) {
	rtr := New()

	for j := 0; j < 9; j++ {
		path := fmt.Sprintf("/ping/:id/%d/pongs", j)
		rtr.Add(path, http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.WriteHeader(http.StatusOK)
		})
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ping/ball/1/pongs", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rtr.ServeHTTP(w, r)
	}
}

func BenchmarkWithPathParamTwentyRequests(b *testing.B) {
	rtr := New()

	for j := 0; j < 9; j++ {
		path := fmt.Sprintf("/ping/:id/%d/pongs", j)
		rtr.Add(path, http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.WriteHeader(http.StatusOK)
		})
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ping/ball/8/pongs", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rtr.ServeHTTP(w, r)
	}
}
