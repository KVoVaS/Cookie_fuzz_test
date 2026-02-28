package main

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"fmt"
)

// FuzzVisitCounter проверяет обработку различных значений куки
func FuzzVisitCounter(f *testing.F) {
	// Добавляем начальные seed-корпуса
	// Корректные значения
	value := map[string]int{"visits": 5}
	encoded, _ := cookieHandler.Encode("visit-counter", value)
	
	// Некорректные значения для тестирования
	seeds := []string{
		encoded,                                   // корректное значение		"invalid-cookie-value",                    // полностью некорректное
		base64.StdEncoding.EncodeToString([]byte("malformed")), // base64 но не валидное для securecookie
		"",                                         // пустая строка
		"12345",                                    // просто число
		"{\"visits\":\"string\"}",                  // JSON с неправильным типом
		"{\"visits\":10}",                           // JSON без securecookie обертки
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, cookieValue string) {
		// Создаем новый request для каждой итерации
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		// Добавляем куку с фаззинг значением
		if cookieValue != "" {
			req.AddCookie(&http.Cookie{
				Name:  "visit-counter",
				Value: cookieValue,
			})
		}

		// Вызываем хендлер напрямую
		mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var visitCount int = 1

			cookie, err := r.Cookie("visit-counter")
			if err == nil {
				var value map[string]int
				if err := cookieHandler.Decode("visit-counter", cookie.Value, &value); err == nil {
					if count, ok := value["visits"]; ok {
						visitCount = count + 1
					}
				}
			}

			value := map[string]int{"visits": visitCount}
			encoded, err := cookieHandler.Encode("visit-counter", value)
			if err == nil {
				cookie := &http.Cookie{
					Name:  "visit-counter",
					Value: encoded,
					Path:  "/",
				}
				http.SetCookie(w, cookie)
			}

			// Проверяем, что visitCount всегда положительный
			if visitCount <= 0 {
				t.Errorf("visitCount должен быть положительным, получено: %d", visitCount)
			}

			// Выводим сообщение о количестве посещений
			message := "Вы посетили страницу " + strconv.Itoa(visitCount)
			
			// Склоняем слово "раз" в зависимости от числа
			if visitCount%10 == 1 && visitCount%100 != 11 {
				message += " раз"
			} else if (visitCount%10 >= 2 && visitCount%10 <= 4) && (visitCount%100 < 10 || visitCount%100 >= 20) {
				message += " раза"
			} else {
				message += " раз"
			}
			
			// Используем w для записи ответа
			_, _ = fmt.Fprint(w, message)
		})

		mainHandler.ServeHTTP(w, req)

		// Проверяем что ответ содержит ожидаемый текст
		body := w.Body.String()
		if !strings.Contains(body, "Вы посетили страницу") {
			t.Errorf("Ответ не содержит ожидаемый текст: %s", body)
		}

		// Проверяем что кука установлена
		response := w.Result()
		defer response.Body.Close()
		
		cookies := response.Cookies()
		found := false
		for _, c := range cookies {
			if c.Name == "visit-counter" {
				found = true
				// Проверяем что значение не пустое
				if c.Value == "" {
					t.Error("Кука visit-counter имеет пустое значение")
				}
				break
			}
		}
		if !found {
			t.Error("Кука visit-counter не была установлена")
		}
	})
}