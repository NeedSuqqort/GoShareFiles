package middleware

import (
	"net/http"
	"context"
	"strings"
)

type BreadCrumbData struct {
	DisplayName string
	URL string
	IsEndOfPath bool
}

type ctx string
const breadcrumbKey ctx = "breadcrumbs"

func BreadcrumbMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		curpath := "/uploads/"
		path := strings.Split(strings.TrimPrefix(request.URL.Path, "/uploads/"),"/")
		breadcrumbs := []BreadCrumbData{
			{DisplayName: "Home", URL: "/uploads/", IsEndOfPath: false},
		}
		for _, folderName := range path {
			if folderName == "" {
				continue
			}
			curpath += folderName + "/"
			breadcrumbs = append(breadcrumbs, BreadCrumbData{DisplayName: folderName, URL: curpath, IsEndOfPath: false})
		}
		breadcrumbs[len(breadcrumbs)-1].IsEndOfPath = true
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