# Запустить фаззинг тест (будет выполняться пока не найдена ошибка или не остановлен)
go test -fuzz=FuzzVisitCounter -fuzztime=30s

# Запустить с определенным временем
go test -fuzz=FuzzVisitCounter -fuzztime=1m

# Запустить все фаззинг тесты
go test -fuzz=. -fuzztime=30s

# Запустить без фаззинга (обычное тестирование)
go test -v