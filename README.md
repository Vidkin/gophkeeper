# gophkeeper

## Начало работы
- проект содержит пример конфигурационного файла для клиента со значениями по умолчанию: cfgclient.yaml.example 
  (необходимо переименовать в cfgclient.yaml для использования)
- проект содержит пример конфигурационного файла для сервера со значениями по умолчанию: cfgserver.yaml.example
  (необходимо переименовать в cfgserver.yaml для использования)
- бинарные данные, загруженные клиентом хранятся на сервере в хранилище MinIO (подключение защищено TLS)
    - при запуске сервера необходимо передать ключи: 
        - -minio-endpoint - адрес сервера MinIO
        - -minio-secret - пользователь пользователя MinIO
        - -minio-id - пользователь MinIO 
- текстовые данные, банковские карты, пары логин/пароль хранятся в БД postgresql в зашифрованном виде
    - при запуске необходимо передать ключи:
        - -d - DSN для подключения к postgresql
        - -db-key - ключ для шифрования логинов и паролей пользователей
- протокол обмена между клиентом и сервером: gRPC (защищён TLS) 
    - при запуске сервера необходимо указать ключи:
        - -crypto-key-private - путь к приватному ключу
        - -crypto-key-public - путь к публичному ключу
- для шифрования JWT при запуске сервера необходимо указать ключ -j
- для расчёта хэша передаваемых данных при запуске сервера необходимо указать ключ -k
- для расчёта хэша передаваемых данных при запуске клиента необходимо указать ключ --hash_key
- шифрование и расшифровку данных из БД осуществляет клиент с помощью ключа --secret_key

### Сборка сервера и клиента + инициализация инфраструктуры со значениями по умолчанию
- обязательно авторизуемся в docker'е:
```
docker login
```
- переходим в корень проекта и выполняем скрипт build_and_start.sh. При этом:
  - автоматически сгенерируются открытый и закрытый ключ;
  - соберётся и запустится контейнер с postgresql и minio;
  - соберётся server и client в корне проекта;
  - будет создан конфиг файл клиента со значениями по умолчанию;
  - запустится сервер со значениями по умолчанию.
- открываем второе окно терминала, переходим в корень проекта и выполняем команды клиента, например:
    - ./client register test test
    - ./client auth test test
    - ./client cards add --owner "Name Surname" --cvv 123 --expire 2024-04-23 --number 13473812 --desc "Test description"
    - ./client cards getAll
    - ./client cards remove --id 1
    - ./client notes add --text "Some text" --desc "Test description"
    - ./client notes getAll
    - ./client notes remove --id 1
    - ./client credentials add --login TestLogin --pass 123 --desc "Test description"
    - ./client credentials getAll
    - ./client credentials remove --id 1
    - ./client files upload --path "/Users/skim/Downloads/Открытый вебинар «Разработка Cloud Native приложений на Go (Введение в Kubernetes)» .mp4" --desc "File description"
    - ./client files getAll --config ./cfgclient.yaml
    - ./client files download --name "Открытый вебинар «Разработка Cloud Native приложений на Go (Введение в Kubernetes)» .mp4" --dir "/Users/skim/Downloads/test"
    - ./client files remove --name "Открытый вебинар «Разработка Cloud Native приложений на Go (Введение в Kubernetes)» .mp4"

### Генерация открытого и закрытого ключа:
Пример команды для генерации открытого и закрытого ключа из корня проекта:
```
go run ./pkg/cert/main.go organization country ./certs/public.crt ./certs/private.key
```

### Запуск контейнера MinIO и PostgreSQL
Выполняем команду: 
```
PRIVATE_KEY_PATH=${путь_к_папке_проекта}/certs/private.key PUBLIC_KEY_PATH=${путь_к_папке_проекта}/certs/public.crt docker-compose up -d
```
Где:

PRIVATE_KEY_PATH - путь к закрытому ключу
PUBLIC_KEY_PATH - путь к открытому ключу

### Сборка сервера
Переходим в каталог cmd/server и выполняем команду:
```
go build -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'"
```

### Сборка клиента
Переходим в каталог cmd/client и выполняем команду:
```
go build -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'"
```

## Запуск сервера

Запускаем сервер, указывая реальные пути к требуемым файлам во флагах. Пример команды:
```
./server -d postgres://postgres:postgres@localhost:5432/postgres -a 127.0.0.1:8080 -crypto-key-public /Users/skim/GolandProjects/gophkeeper/certs/public.crt -crypto-key-private /Users/skim/GolandProjects/gophkeeper/certs/private.key -db-key strongDBKey2Ks5nM2J5JaI59PPEhL1x -j JWTKey -minio-endpoint 127.0.0.1:9000 -minio-secret minioadmin -minio-id minioadmin -k defaultHashKey
```

## Пример команд клиента

### Регистрация пользователя с логином test и паролем test
```
./client register test test --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

### Авторизация пользователя с логином test и паролем test
```
./client auth test test --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

### Банковские карты
#### Добавление новой банковской карты
```
./client cards add --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x --owner "Name Surname" --cvv 123 --expire 2024-04-23 --number 13473812 --desc "Test description"
```

#### Показать все банковские карты
```
./client cards getAll --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

#### Показать банковскую карту по id
```
./client cards get --id 1 --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

#### Удалить банковскую карту по id
```
./client cards remove --id 1 --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

### Текстовые данные

#### Добавление новых текстовых данных
```
./client notes add --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x --text "Some text" --desc "Test description"
```

### Показать все текстовые данные
```
./client notes getAll --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

### Показать текстовые данные по id
```
./client notes get --id 1 --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

### Удалить текстовы данные по id
```
./client notes remove --id 1 --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

### Пары логин-пароль

#### Добавление новой пары логин-пароль
```
./client credentials add --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x --login TestLogin --pass 123 --desc "Test description"
```

#### Показать все пары логин-пароль
```
./client credentials getAll --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

#### Показать пару логин-пароль по id
```
./client credentials get --id 1 --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

#### Удалить пару логин-пароль по id
```
./client credentials remove --id 1 --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

### Файлы

#### Upload
```
./client files upload --path "/Users/skim/Downloads/Открытый вебинар «Разработка Cloud Native приложений на Go (Введение в Kubernetes)» .mp4" --config ./cfgclient.yaml --desc "File description" --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

#### Показать все файлы
```
./client files getAll --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

#### Download
```
./client files download --name "Открытый вебинар «Разработка Cloud Native приложений на Go (Введение в Kubernetes)» .mp4" --dir "/Users/skim/Downloads/test" --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```

#### Удалить файл по имени
```
./client files remove --name "Открытый вебинар «Разработка Cloud Native приложений на Go (Введение в Kubernetes)» .mp4" --config ./cfgclient.yaml --hash_key defaultHashKey --secret_key strongDBKey2Ks5nM2J5JaI59PPEhL1x
```