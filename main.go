package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Переменная для хранения счетчика
		var visitCount int = 1

		// Чтение существующей куки
		cookie, err := r.Cookie("visit-counter")
		if err == nil {
			// Декодируем значение куки
			var value map[string]int
			if err := cookieHandler.Decode("visit-counter", cookie.Value, &value); err == nil {
				// Увеличиваем счетчик на 1
				if count, ok := value["visits"]; ok {
					visitCount = count + 1
				}
			}
		}

		// Создаем или обновляем куку с новым значением счетчика
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

		fmt.Fprint(w, message)
	})

	http.ListenAndServe(":8080", nil)
}
