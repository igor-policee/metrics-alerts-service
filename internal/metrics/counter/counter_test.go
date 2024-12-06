// internal/metrics/counter/counter_test.go

// TODO: change language

package counter

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateCounterHandler(t *testing.T) {
	// Инициализация
	req, err := http.NewRequest("POST", "/update/counter/testCounter/10", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	// Вызов обработчика
	UpdateCounterHandler(rr, req)

	// Проверка статуса ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус %v, получен %v", http.StatusOK, status)
	}

	// Проверка содержимого ответа
	expected := "Counter metric 'testCounter' updated to 10 successfully."
	if rr.Body.String() != expected {
		t.Errorf("Ожидался ответ '%s', получен '%s'", expected, rr.Body.String())
	}
}

func TestGetCounterHandler(t *testing.T) {
	// Предварительная установка значения счетчика
	mu.Lock()
	counterMetrics["testCounter"] = 15
	mu.Unlock()

	// Создание запроса
	req, err := http.NewRequest("GET", "/value/counter/testCounter", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Вызов обработчика
	GetCounterHandler(rr, req)

	// Проверка статуса ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус %v, получен %v", http.StatusOK, status)
	}

	// Проверка содержимого ответа
	expected := "15"
	if rr.Body.String() != expected {
		t.Errorf("Ожидался ответ '%s', получен '%s'", expected, rr.Body.String())
	}
}

func TestGetCounterHandler_NotFound(t *testing.T) {
	// Создание запроса для несуществующего счетчика
	req, err := http.NewRequest("GET", "/value/counter/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Вызов обработчика
	GetCounterHandler(rr, req)

	// Проверка статуса ответа
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Ожидался статус %v, получен %v", http.StatusNotFound, status)
	}
}
