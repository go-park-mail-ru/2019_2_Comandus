#### Добавить пользователя:
`$: http http://localhost:8080/users name=dasha password=1234 email=dasha@ya.ru`
#### Создать сессию: 
по умолчанию заходим как фрилансер   
`$: http -v --session=user POST http://localhost:8080/sessions password=1234 email=dasha@ya.ru`
#### Задать тип пользователя в куки
`$: http -v --session=user http://localhost:8080/private/setusertype user_type=customer`
или  
`$: http -v --session=user http://localhost:8080/private/setusertype user_type=freelancer`
#### Попробовать запустить приватную(только для авторизованных) функцию:
`$: http -v --session=user http://localhost:8080/private/profile`
#### Удалить сессию:
`$: http -v --session=user http://localhost:8080/private/logout`

###### Завести БД: (пока не нужно)
`$: psql restapi_dev`  
`restapi_dev=> ALTER USER d PASSWORD '1234';`

#### Сборка
`$: make build # port 8080 as default` 

#### Запуск сервера
`$: ./apiserver`

#### Запуск тестов
`$: make test`