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

// self work func
func TestCafeWhenFail(t *testing.T) {
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

	// /cafe — должно возвращаться "unknown city" с кодом http.StatusBadRequest;
	// /cafe?city=omsk — должно возвращаться "unknown city" с кодом http.StatusBadRequest;
	// /cafe?city=tula&count=na — должно возвращаться "incorrect count" с кодом http.StatusBadRequest;

}

// TestCafeCount() — проверяет работу сервера при разных значениях параметра count;
func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	// slices fo loops
	citiesSlice := []string{"moscow", "tula"}
	countSlice := []int{0, 1, 2, 100}
	for _, city := range citiesSlice {
		for _, count := range countSlice {
			requestString := fmt.Sprintf("/cafe?city=%s&count=%d", city, count)
			response := httptest.NewRecorder()
			req := httptest.NewRequest("GET", requestString, nil)
			handler.ServeHTTP(response, req)
			// fetch response and convert to string
			responseString := response.Body.String()
			// string to slice, separator ","
			responseStringToSlice := strings.Split(responseString, ",")
			// DEBUG PRINT
			// fmt.Println(len(responseString), responseStringToSlice, len(responseStringToSlice), count)
			if count == 0 || len(responseString) == 0 {
				assert.Equal(t, len(responseString), count) // len(responseStringToSlice) with "" - server is returning an empty string, which when split by comma, results in a slice containing one empty string.
			} else if count == 1 {
				assert.Equal(t, len(responseStringToSlice), count)
			} else if count == 2 {
				assert.Equal(t, len(responseStringToSlice), count)
			} else if count == 100 {
				// var wantCount int
				if city == "moscow" {
					wantCount := 5
					assert.Equal(t, len(cafeList[city]), wantCount)
				}
				if city == "tula" {
					wantCount := 3
					assert.Equal(t, len(cafeList[city]), wantCount)
				}
			}
		}
	}
}

// TestCafeSearch() — проверяет результат поиска кафе по указанной подстроке в параметре search.
func TestCafeSearch(t *testing.T) {
	requests := []struct {
		search    string // передаваемое значение search
		wantCount int    // ожидаемое количество кафе в ответе
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	handler := http.HandlerFunc(mainHandle)
	for _, v := range requests {
		requestString := fmt.Sprintf("/cafe?city=moscow&search=%s", v.search)
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", requestString, nil)
		handler.ServeHTTP(response, req)
		// С помощью require.Equal() следует проверить, что запрос успешно обработан,
		// Проверяем, что запрос успешно обработан
		require.Equal(t, http.StatusOK, response.Code)
		// fetch response and convert to string
		responseString := strings.ToLower(response.Body.String()) // no need responseString = strings.ToLower(responseString), we can wrap "response.Body.String()"
		if strings.Contains(responseString, v.search) {
			// https://stackoverflow.com/questions/21417987/how-to-get-number-of-results-in-strings-contains
			// strings.Count(responseString, v.search) - counts number of entries of subString( v.search)
			assert.Equal(t, strings.Count(responseString, v.search), v.wantCount)
		}
	}
}
