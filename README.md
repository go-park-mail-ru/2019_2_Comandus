#### Добавить пользователя:
`$: http http://localhost:8080/signup name=dasha password=1234 email=dasha@ya.ru`

`$ curl -XPOST http://127.0.0.1:8080/users --data '{"name" : "dasha" , "password" : "1234" , "email" : "dasha@ya.ru"}'`
#### Создать сессию: 
по умолчанию заходим как фрилансер   
`$: http -v --session=user POST http://localhost:8080/login password=1234 email=dasha@ya.ru`

`$ 'curl -XPOST -v -c cookie.txt  http://127.0.0.1:8080/sessions --data '{"password" : "1234" , "email" : "dasha@ya.ru"}'`
#### Задать тип пользователя в куки
`$: http -v --session=user http://localhost:8080/private/setusertype user_type=customer` 
или 
`$: http -v --session=user http://localhost:8080/private/setusertype user_type=freelancer`

#### Попробовать запустить приватную(только для авторизованных) функцию:
`$: http -v --session=user http://localhost:8080/private/profile`

`$: curl -b cookie.txt -v http://127.0.0.1:8080/private/profile`

#### Добавить изменения в профиль 
`$: curl -b cookie.txt -XPOST -v http://127.0.0.1:8080/private/profile/edit
 --data '{
    "first_name" : "Dima",
    "second_name" : "Andronov" ,
    "contact_info": {
        "country" : "Russia" ,
         "city" : "Moscow" ,
          "phone_number" : "89870720609"}}'`

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