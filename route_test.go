package router

import (
	"fmt"
	"github.com/shyamz-22/router/assert"
	"github.com/shyamz-22/router/fixture"

	"net/http"
	"net/http/httptest"
	"testing"
)

var pong = []byte("Pong!")
var pongBack = []byte("Pong back!")

func TestRoute(t *testing.T) {
	t.Parallel()
	t.Run("adds a static route", func(t *testing.T) {
		rtr := New()
		rtr.Add("/ping", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pong)
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/ping", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusOK, "Pong!")
	})

	t.Run("adds a static route with more paths", func(t *testing.T) {
		rtr := New()
		rtr.Add("/ping/another", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pongBack)
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/ping/another", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusOK, "Pong back!")
	})

	t.Run("adds a static route with file extension", func(t *testing.T) {
		rtr := New()
		rtr.Add("/articles/", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pong)
		})

		rtr.Add("/articles/go_command.html", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pongBack)
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/articles/go_command.html", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusOK, "Pong back!")
	})

	t.Run("adds a static route with leading slash", func(t *testing.T) {
		rtr := New()
		rtr.Add("/articles/", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pong)
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/articles/", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusOK, "Pong!")
	})
}

func TestRouteWithNotFoundPaths(t *testing.T) {
	t.Parallel()
	t.Run("returns 404 for not found path", func(t *testing.T) {
		rtr := New()
		rtr.Add("/ping/another", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pongBack)
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/ping", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithStatus(t, w, http.StatusNotFound)
	})

	t.Run("returns 404 for not found path with similar base paths", func(t *testing.T) {
		rtr := New()
		rtr.Add("/ping", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pongBack)
		})
		rtr.Add("/ping/another", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pongBack)
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/ping/another/yet", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithStatus(t, w, http.StatusNotFound)
	})

}

func TestRouteWithMethodNotAllowed(t *testing.T) {
	t.Parallel()
	rtr := New()
	rtr.Add("/pings", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
		w.Write(pong)
	})

	rtr.Add("/pings/:id", http.MethodPost, func(w http.ResponseWriter, r *http.Request, params PathParams) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Pong post!"))
	})

	t.Run("static route", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/pings", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithStatus(t, w, http.StatusMethodNotAllowed)
	})

	t.Run("route with path params", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodDelete, "/pings/another", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithStatus(t, w, http.StatusMethodNotAllowed)
	})
}

func TestRouteWithPathParams(t *testing.T) {
	t.Parallel()
	t.Run("returns 200 for a route with path param", func(t *testing.T) {
		rtr := New()
		rtr.Add("/ping/:id", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			id := params.ByName("id")
			w.Write([]byte(id))
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/ping/Pong", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusOK, "Pong")
	})

	t.Run("returns 200 for a route with multiple path param", func(t *testing.T) {
		rtr := New()
		rtr.Add("/ping/:pingId/pong/:pongId", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pongBack)
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/ping/1/pong/2", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusOK, "Pong back!")
	})

	t.Run("returns 200 for a route path param in the middle", func(t *testing.T) {
		rtr := New()
		rtr.Add("/ping/:pingId/yet/pongs", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write([]byte("4s"))
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/ping/some-id/yet/pongs", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusOK, "4s")
	})

	t.Run("returns 404 for not found path with similar base paths", func(t *testing.T) {
		rtr := New()
		rtr.Add("/ping", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pongBack)
		})
		rtr.Add("/ping/:pingId/pong/:pongId", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
			w.Write(pongBack)
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/ping/another/pong/yet/ping/pong", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithStatus(t, w, http.StatusNotFound)
	})
}

func TestRouteWithMultipleMethods(t *testing.T) {
	t.Parallel()
	rtr := New()
	rtr.Add("/pings", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
		w.Write(pong)
	})

	rtr.Add("/pings/:id", http.MethodPost, func(w http.ResponseWriter, r *http.Request, params PathParams) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Pong post!"))
	})

	rtr.Add("/pings/:id", http.MethodDelete, func(w http.ResponseWriter, r *http.Request, params PathParams) {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Pong delete!"))
	})

	t.Run("returns 200 for a route with Get", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/pings", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusOK, "Pong!")
	})

	t.Run("returns 201 for a route with Post", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/pings/1", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusCreated, "Pong post!")
	})

	t.Run("returns 204 for a route with Delete", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodDelete, "/pings/1", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusNoContent, "Pong delete!")
	})

}

func TestRouteWithUncleanPath(t *testing.T) {
	t.Parallel()
	rtr := New()
	rtr.Add("/", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
		w.Write([]byte("Index"))
	})

	rtr.Add("/pings/:id/pongs/:pongId", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params PathParams) {
		w.Write([]byte(params.ByName("pongId")))
	})

	t.Run("with empty path param", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/pings//pongs/pong", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusOK, "pong")
	})

	t.Run("trailing slash at end is not ignored", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/pings/ping/pongs//", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithStatus(t, w, http.StatusNotFound)
	})

	t.Run("route with just /", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		rtr.ServeHTTP(w, r)

		assert.ResponseWithBody(t, w, http.StatusOK, "Index")
	})
}

func TestRouter_ServeHTTP(t *testing.T) {
	rtr := New()
	for _, route := range fixture.Routes {
		rtr.Add(route.Path, route.Method, func(writer http.ResponseWriter, request *http.Request, params PathParams) {
			writer.WriteHeader(http.StatusOK)
		})
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	u := r.URL
	rq := u.RawQuery

	for _, route := range fixture.RoutesWithPathValues {
		r.Method = route.Method
		r.RequestURI = route.Path
		u.Path = route.Path
		u.RawQuery = rq
		fmt.Println(r.RequestURI)
		rtr.ServeHTTP(w, r)
		assert.ResponseWithStatus(t, w, http.StatusOK)
	}

}
