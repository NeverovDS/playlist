Тестовое задание для поступления в GoCloudCamp

Для запуска:

`cd ./playlist/cmd/playlist`
`go run main.go`

- API port `localhost:8000`

Тело запросов: { "name" : "SongName", "duration" : "12345678" }

Эндпоинты:

- GET /playlist/play
- GET /playlist/pause
- GET /playlist/next
- GET /playlist/prev
- POST /playlist/create/song
- GET /playlist/songs - список песен
- GET /playlist/song/:id - показать песню
- PUT /playlist/song/:id - обновить песню
- DELETE /playlist/song/:id - удалить песню