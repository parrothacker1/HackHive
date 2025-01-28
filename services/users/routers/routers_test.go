package routers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
  os.Setenv("GO_ENV","test")
  m.Run()
}

func TestRouter(t *testing.T) {
  router_test := NewRouter()
  test_handler := func (w http.ResponseWriter,r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("1"))
  }
  middleware := func(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      if (r.Header.Get("Deny") == "1") {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("0"))
        return
      }
      next.ServeHTTP(w,r)
    }
  }
  router_test.Handle(http.MethodGet,"/",test_handler,middleware)
  t.Run("Testing with normal request",func(t *testing.T) {
    rr := httptest.NewRecorder()
    req,err := http.NewRequest(http.MethodGet,"/",nil)
    if err != nil {
      t.Fatal("Failed to create request for testing")
    }
    router_test.ServeHTTP(rr,req)
    require.Equal(t,http.StatusOK,rr.Code,"The Router is not working.The return value should be 200")
  })
  t.Run("Testing with wrong method",func(t *testing.T) {
    rr := httptest.NewRecorder()
    req_post,err := http.NewRequest(http.MethodPost,"/",nil)
    if err != nil {
      t.Fatal("Failed to create request for testing")
    }
    router_test.ServeHTTP(rr,req_post)
    require.Equal(t,http.StatusNotFound,rr.Code,"The Router is not working.The return value should be 404")
  })
  t.Run("Testing with middleware",func(t *testing.T) {
    req,err := http.NewRequest(http.MethodGet,"/",nil)
    if err != nil {
      t.Fatal("Failed to create request for testing")
    }
    req.Header.Add("Deny","1")
    rr := httptest.NewRecorder()
    router_test.ServeHTTP(rr,req)
    require.Equal(t,http.StatusBadRequest,rr.Code,"The Router middleware is not working.The return value should be 400")
  })
  t.Run("Testing with wrong path",func(t *testing.T) {
    req,err := http.NewRequest(http.MethodGet,"/wrongpath",nil)
    if err != nil {
      t.Fatal("Failed to create request for testing")
    }
    rr := httptest.NewRecorder()
    router_test.ServeHTTP(rr,req)
    require.Equal(t,http.StatusNotFound,rr.Code,"The Router is not working.The return value should be 404")
  })
}
