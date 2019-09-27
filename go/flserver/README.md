Добавить пользователя:
http http://localhost:8080/users name=dasha password=1234
Создать сессию:
http -v --session=user POST http://localhost:8080/sessions name=dasha password=1234
Попробовать запустить приватную(только для авторизованных) функцию:
http -v --session=user http://localhost:8080/private/profile


Завести БД:
restapi_dev=> ALTER USER d PASSWORD '1234';
