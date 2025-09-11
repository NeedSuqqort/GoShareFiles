package middleware

import (
	"net/http"
	"context"
	"strings"
)

type BreadCrumbData struct {
	DisplayName string
	URL string
}

type ctx string
const breadcrumbKey ctx = "breadcrumbs"

func BreadcrumbMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		curpath := "/uploads/"
		breadcrumbs := []BreadCrumbData{
			{DisplayName: "Home", URL: "/uploads/"},
		}
		path := strings.Split(strings.TrimPrefix(request.URL.Path, "/uploads/"),"/")
		for _, folderName := range path {
			curpath += folderName + "/"
			breadcrumbs = append(breadcrumbs, BreadCrumbData{DisplayName: folderName, URL: curpath})
		}
		ctx := context.WithValue(request.Context(), breadcrumbKey, breadcrumbs)
        next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetBreadcrumbs(request *http.Request) []BreadCrumbData {
	if breadcrumbs, ok := request.Context().Value(breadcrumbKey).([]BreadCrumbData); ok {
		return breadcrumbs
	} else {
		return nil
	}
}