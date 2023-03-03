package recoverer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestPanic(t *testing.T) {

	panicHandler := func(w http.ResponseWriter, r *http.Request) {
		panic("At the Disco")
	}

	l, _ := zap.NewProduction()
	s := httptest.NewServer(Recoverer(http.HandlerFunc(panicHandler), l))

	request, _ := http.NewRequest("GET", s.URL, nil)
	http.DefaultClient.Do(request)
}
