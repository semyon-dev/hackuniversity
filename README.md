# hackuniversity 2020
team's name: Time walkers

## Technology stack
* Go v1.14
* Clickhouse
* PostgresSQL
* OPC UA Simulator
* Gin
* Gorilla websockets
* JS OPC UA

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
API отвечает за изменение min и max параметров (критические параметры)

## opc
`npm install node-opcua` \
`node opc.js` \
opc - симулятор opc ua server

## Scheme
![](https://github.com/semyon-dev/hackuniversity/blob/master/scheme.png) 

## LICENSE
MIT