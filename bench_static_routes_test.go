package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkWithOneRequest(b *testing.B) {

	rtr := New()
	rtr.Add("/ping", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
		w.Write(pong)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ping", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rtr.ServeHTTP(w, r)
	}
}

func BenchmarkWithTenRequests(b *testing.B) {
	rtr := New()

	for j := 0; j < 9; j++ {
		path := fmt.Sprintf("/ping/%d", j)
		rtr.Add(path, http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pong)
		})
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ping/1", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rtr.ServeHTTP(w, r)
	}
}

func BenchmarkWithTwentyRequests(b *testing.B) {
	rtr := New()

	for j := 0; j < 9; j++ {
		path := fmt.Sprintf("/ping/%d", j)
		rtr.Add(path, http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pong)
		})
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ping/8", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rtr.ServeHTTP(w, r)
	}
}
