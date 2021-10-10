# url_shortener
URL shortener service

endpoints:
POST:/links - добавление ссылки
DELETE:/links/{id} - удаление ссылки
GET:/links?phrase=goog - получение списка ссылок или по ключевому слову в длинной версии ссылки или всего списка целиком, если не указывать query параметр