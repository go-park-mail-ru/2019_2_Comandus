#### Генерация mock-файлов для структур Repository и Usecase
`$: ./internal/app/mocks/generate_mocks.sh` 
##### Mock-файлы для структур типа Usecase сразу создаются в папках internal/app/{название модели}/delivery/http/test
##### Mock-файлы для структур типа Repository создаются в папке internal/app/test, т.к. структуры типа Usecase могут использовать несколько репозиториев

