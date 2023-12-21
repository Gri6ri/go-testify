package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return

	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}
func TestMainHandlerIsRequestOkAndNotEmpty(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=3&city=moscow", nil)
	if err != nil {
		fmt.Println("Ошибка формирования запроса:", err)
		return
	}
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	require.Equal(t, http.StatusOK, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Body)
}
func TestMainHandlerIfCityIsNotSupported(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=1&city=kazan", nil)
	if err != nil {
		fmt.Println("Ошибка формирования запроса:", err)
		return
	}
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	expectedCity := "moscow"
	actualCity := req.URL.Query().Get("city")
	assert.Equal(t, expectedCity, actualCity)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	expectedError := "wrong count value"
	actualError := responseRecorder.Body.String()
	assert.Equal(t, expectedError, actualError)
}
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req, err := http.NewRequest("GET", "/cafe?count=8&city=moscow", nil)
	if err != nil {
		fmt.Println("Ошибка формирования запроса:", err)
		return
	}
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	countStr := req.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		fmt.Println("Ошибка преобразования переменной countStr в int", err)
		return
	}
	assert.Greater(t, count, totalCount)
	expectedAnswer := "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент"
	actualAnswer := responseRecorder.Body.String()
	assert.Equal(t, expectedAnswer, actualAnswer)

}
