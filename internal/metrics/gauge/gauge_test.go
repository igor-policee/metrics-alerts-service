// internal/metrics/gauge/gauge_test.go

// TODO: change language

package gauge

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateGaugeHandler(t *testing.T) {
	// Инициализация
	req, err := http.NewRequest("POST", "/update/gauge/testGauge/75.5", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	// Вызов обработчика
	UpdateGaugeHandler(rr, req)

	// Проверка статуса ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус %v, получен %v", http.StatusOK, status)
	}

	// Проверка содержимого ответа
	expected := "Gauge metric 'testGauge' set to 75.50 successfully."
	if rr.Body.String() != expected {
		t.Errorf("Ожидался ответ '%s', получен '%s'", expected, rr.Body.String())
	}
}

func TestGetGaugeHandler(t *testing.T) {
	// Предварительная установка значения gauge
	mu.Lock()
	gaugeMetrics["testGauge"] = 80.25
	mu.Unlock()

	// Создание запроса
	req, err := http.NewRequest("GET", "/value/gauge/testGauge", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Вызов обработчика
	GetGaugeHandler(rr, req)

	// Проверка статуса ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус %v, получен %v", http.StatusOK, status)
	}

	// Проверка содержимого ответа
	expected := "80.25"
	if rr.Body.String() != expected {
		t.Errorf("Ожидался ответ '%s', получен '%s'", expected, rr.Body.String())
	}
}

func TestGetGaugeHandler_NotFound(t *testing.T) {
	// Создание запроса для несуществующего gauge
	req, err := http.NewRequest("GET", "/value/gauge/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Вызов обработчика
	GetGaugeHandler(rr, req)

	// Проверка статуса ответа
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Ожидался статус %v, получен %v", http.StatusNotFound, status)
	}
}
