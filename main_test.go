package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"

	request := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, min(len(cafeList[city]), 100)},
	}

	for _, v := range request {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?count=%d&city=%s", v.count, city), nil)

		handler.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)

		if res.Body.String() == "" {
			fmt.Printf("expected %d, got %d\n", v.want, len(res.Body.String()))
			assert.Equal(t, v.want, len(res.Body.String()))
			continue
		}
		fmt.Printf("expected %d, got %d\n", v.want, len(strings.Split(res.Body.String(), ",")))
		assert.Equal(t, v.want, len(strings.Split(res.Body.String(), ",")))

	}

}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=%s&search=%s", city, v.search), nil)

		handler.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)

		var foundedCafeList []string
		for _, name := range strings.Split(strings.TrimSpace(res.Body.String()), ",") {

			if strings.Contains(strings.ToLower(name), strings.ToLower(v.search)) {
				foundedCafeList = append(foundedCafeList, name)
			}
		}
		fmt.Printf("expected count cafe %d, got count %d\n", v.wantCount, len(foundedCafeList))
		assert.Equal(t, v.wantCount, len(foundedCafeList))

	}
}
