# serialization_estimator

Приложение, представляющее из себя udp-прокси и несколько udp-серверов,
позволяющее тестировать эффективность различных форматов сериализации.

#### Поддерживаемые форматы

native (encoding/gob), json, xml, msgpack, yaml, protobuf и avro

### Автоматическая сборка и запуск
Достаточно в корне репозитория запустить
```bash
$ docker compose up
```

### Ручная сборка
Чтобы собрать все образы, нужно запустить исполнить скрипт [build_images.sh](https://github.com/fdr896/serialization_estimator/blob/main/build_images.sh)
Запушить образы в докер-хаб возможно с помощью скрипта [update_images.sh](https://github.com/fdr896/serialization_estimator/blob/main/update_images.sh)

## Работа сервиса
Запущенный через `docker compose сервис, ожидает udp-соединений на порту 2000.

Сервис принимает запросы двух видов:
- `{"method": "get_result", "param":"<serialization_type>"}` (`<serialization_type>` нужно заменить на один из вышеперечисленных способов сериализации, например, `{"method": "get_result", "param":"protobuf"}`)
- `{"method": "get_result_all"}` -- прокси рассылает запрос ко всем estimator'ам на адрес мультикаст группы, после чего пересылает все ответы пользователю

### Тестирование при помощи `netcat`
Пример взаимодействия с сервисом (предполагается, что сервис был запущен) (запросы нужно вписывать вручную текстом)
```bash
$ cat | netcat -u 127.0.0.1 2000
{"method":"get_result", "param":"protobuf"}
protobuf - 189b - 2741ns - 3485ns
{"method":"get_result", "param":"native"}    
native - 285b - 5617ns - 34ns
{"method":"get_result_all"}
avro - 168b - 1252ns - 2510ns
msgpack - 182b - 2698ns - 3845ns
native - 285b - 9138ns - 53ns
protobuf - 189b - 5717ns - 5347ns
json - 227b - 3364ns - 8240ns
xml - 277b - 8896ns - 22969ns
yaml - 220b - 20125ns - 26754ns
{"method":"get_result", "param":"xml"}
xml - 277b - 8085ns - 23504ns
```
