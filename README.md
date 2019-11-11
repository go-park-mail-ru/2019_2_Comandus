# 2019_2_Comandus

### How to Run | Build

**Run Mode**

- `make build`
- `./apiserver`

**Test Mode**

- `./internal/app/test/generate_mocks.sh`  
- `make test`

## Проект

fl.ru ( ставки на проект )

**Features:** 

- Регистрация как фрилансер и как заказчик (по возможности OAuth 
- Разместить заказ (категория, название, описание, бюджет)
- Возможность откликаться на заказ фрилансерам (ставки на проект)
- Разместить вакансию
- Отзывы и оценки (и заказу и фрилансерам)
- Хороший поиск с фильтром

**REST API**

https://swagger.nozim.now.sh/

**Макеты UI**

- [Папка с макетами](docs/wireframes)

- [Ninjamock (старые версии)](https://ninjastorage.blob.core.windows.net/html/SMFDQFx/52e2914d-427c-06f1-ebb8-e593fdbce622.html)

**Kanban Board**

https://trello.com/comandus

### Команда

- [Дарья Ефимова](https://github.com/efimovad)
- [Дмитрий Андронов](https://github.com/Andronovdima)
- [Нозим Юнусов](https://github.com/nozimy)
- [Александр Косенков](https://github.com/SoulPhazed)

Ментор: [Джахонгир Тулфоров](https://github.com/bin-umar)

### Frontend

https://github.com/frontend-park-mail-ru/2019_2_Comandus
### Backend

https://github.com/go-park-mail-ru/2019_2_Comandus

### Примеры запросов на сервер
#### Добавить пользователя:
`$: http http://localhost:8080/signup username=dasha password=1234 email=dasha@ya.ru`

`$ curl -XPOST http://127.0.0.1:8080/users --data '{"name" : "dasha1" , "password" : "1234" , "email" : "dasha1@ya.ru"}'`
#### Создать сессию: 
по умолчанию заходим как фрилансер   
`$: http -v --session=user POST http://localhost:8080/login password=1234 email=dasha@ya.ru`

`$ 'curl -XPOST -v -c cookie.txt  http://127.0.0.1:8080/sessions --data '{"password" : "1234" , "email" : "dasha@ya.ru"}'`
#### Задать тип пользователя в куки
`$: http -v --session=user http://localhost:8080/private/setusertype user_type=client` 
или 
`$: http -v --session=user http://localhost:8080/private/setusertype user_type=freelancer`

#### Попробовать запустить приватную(только для авторизованных) функцию:
`$: http -v --session=user http://localhost:8080/private/account`

`$: curl -b cookie.txt -v http://127.0.0.1:8080/private/account`

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
