Добавить пользователя:
http http://localhost:8080/users name=dasha password=1234 email=dasha@ya.ru
Создать сессию:
http -v --session=user POST http://localhost:8080/sessions password=1234 email=dasha@ya.ru
Попробовать запустить приватную(только для авторизованных) функцию:
http -v --session=user http://localhost:8080/private/profile
Удалить сессию:
http -v --session=user http://localhost:8080/private/logout

Завести БД: (пока не нужно)
restapi_dev=> ALTER USER d PASSWORD '1234';

есть Makefile поэтому:
make build
./apiserver

запустить тесты
make test