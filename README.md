# url_shortener
URL shortener service

Сервис контейнеризирован. Сделан фронтенд на шаблонах Go из-за чего кстати не удалось реализовать потоковую передачу данных для отрисовки. В качестве хранилища используется БД Postgres.

Запуск:

docker-compose build

docker-compose up

После чего сервис будет доступен из браузера по ссылке:

http://localhost:8080 или http://your_docker_host_address:8080

P.S. В связи с нехваткой времени пришлось использовать знакомый мне gorilla/mux вместо go-chi.