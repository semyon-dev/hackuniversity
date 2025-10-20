# hackuniversity 2020

Team: time walkers (origin dev), task by BIOCAD

# Stack
* Go v1.14
* Clickhouse
* OPC UA Simulator
* PostgresSQL
* Gin
* Gorilla websockets
* JS OPC UA

## Other parts

* [Frontend (React)](https://github.com/dchudik/hackuniversity2020)
* [Mobile (Android, kotlin)](https://github.com/GranPecador/BiocadProject)

# Микросервисы
## pusher
`go run pusher/main.go` \
Pusher получает данные с OPC server, записывает в clickhouse и отдает их по вебсокетам другим микросервисам

## checkerr
`go run checkerr/main.go` \
checkerr получает данные с pusher (по websocket) и проверят данные по критическим параметрам, \
в случае нахождения превышений - сохраняет в журнал ошибок (postgres) и  \
отправляет ошибки пользователям через telegram бота

## api
`go run api/main.go` \
HTTP API отвечает за получение/изменение min и max параметров (критические параметры). \
А также за аналитику данных из разных баз данных для клиентов

## unloader
`go run unloader/main.go` \
Этот микросервис отвечает за разгрузку данных на клиенты от микросервиса pusher

## opc
`npm install node-opcua` \
`node opc.js` \
opc - симулятор opc ua server который генерирует разные данные типа float каждую секунду

# Схема
![](https://github.com/semyon-dev/hackuniversity/blob/master/scheme.png) 

# LICENSE
[MIT](https://github.com/semyon-dev/hackuniversity/blob/master/LICENSE)
