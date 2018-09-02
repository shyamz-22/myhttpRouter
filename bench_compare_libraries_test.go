package router

import (
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/shyamz-22/router/fixture"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkGithubV3(b *testing.B) {
	rtr := New()
	for _, route := range fixture.Routes {
		rtr.Add(route.Path, route.Method, func(writer http.ResponseWriter, request *http.Request, params PathParams) {
			writer.WriteHeader(http.StatusOK)
		})
	}

	benchRoutes(b, rtr, fixture.RoutesWithPathValues)
}

func BenchmarkGithubV3_hp(b *testing.B) {
	rtr := httprouter.New()
	for _, route := range fixture.Routes {
		rtr.Handle(route.Method, route.Path, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
			writer.WriteHeader(http.StatusOK)
		})
	}

	benchRoutes(b, rtr, fixture.RoutesWithPathValues)
}

func BenchmarkGithubV3_mux(b *testing.B) {
	rtr := mux.NewRouter()
	for _, route := range fixture.MuxRoutes {
		rtr.HandleFunc(route.Path, func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
		}).Methods(route.Method)
	}

	benchRoutes(b, rtr, fixture.RoutesWithPathValues)
}

func benchRoutes(b *testing.B, router http.Handler, routes []fixture.Route) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	u := r.URL
	rq := u.RawQuery

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, route := range routes {
			r.Method = route.Method
			r.RequestURI = route.Path
			u.Path = route.Path
			u.RawQuery = rq
			router.ServeHTTP(w, r)
		}
	}
}
