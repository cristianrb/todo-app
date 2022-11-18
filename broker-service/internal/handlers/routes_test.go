package handlers

//
//import (
//	mocks "authentication/internal/mocks/services"
//	"github.com/go-chi/chi/v5"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"testing"
//)
//
//func Test_routes_exist(t *testing.T) {
//	repo := mocks.NewUsersServiceMocks()
//	handlerConfig := NewHandlerConfig(nil, nil)
//	testRoutes := handlerConfig.Routes()
//	chiRoutes := testRoutes.(chi.Router)
//
//	routes := []string{"/pushItemToQueue"}
//
//	for _, route := range routes {
//		assert.Equal(t, true, routeExists(chiRoutes, route))
//	}
//}
//
//func routeExists(routes chi.Router, route string) bool {
//	found := false
//
//	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
//		if route == foundRoute {
//			found = true
//		}
//		return nil
//	})
//
//	return found
//}
