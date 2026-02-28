# Запустить фаззинг тест
go test -fuzz=FuzzVisitCounter -fuzztime=30s

# Запустить без фаззинга (обычное тестирование)
go test -v