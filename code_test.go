package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func PrepareResponse(url string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest("GET", url, nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	return responseRecorder, nil
}
func TestMainHandlerIsRequestOkAndNotEmpty(t *testing.T) {
	responseRecorder, err := PrepareResponse("/cafe?count=3&city=moscow")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Body)

}
func TestMainHandlerIfCityIsNotSupported(t *testing.T) {
	responseRecorder, err := PrepareResponse("/cafe?count=1&city=kazan")
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	// в задании требуется, чтобы сервис возвращал ошибку "wrong count value".
	// Опечатка, на мой взгляд. Иначе тест не будет проходить.
	expectedError := "wrong city value"
	actualError := responseRecorder.Body.String()
	assert.Equal(t, expectedError, actualError)
}
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	responseRecorder, err := PrepareResponse("/cafe?count=8&city=moscow")
	require.NoError(t, err)
	expectedCount := 4
	actualAnswerSlice := strings.Split(responseRecorder.Body.String(), ",")
	assert.Len(t, actualAnswerSlice, expectedCount)
	expectedAnswer := "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент"
	actualAnswerString := responseRecorder.Body.String()
	assert.Equal(t, expectedAnswer, actualAnswerString)

}
