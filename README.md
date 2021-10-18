# url_shortener
URL shortener service

endpoints:
POST:/links - добавление ссылки
DELETE:/links/{id} - удаление ссылки
GET:/links?phrase=goog - получение списка ссылок или по ключевому слову в длинной версии ссылки или всего списка целиком, если не указывать query параметр

В качестве роутера для курсового проекта решил остановить свой выбор на go-chi.
Раньше я работал с gorilla/mux, и он меня полностью устраивает, однако хочется попрактиковаться в работе с go-chi, к тому же при его использовании код получается немного компактнее.
Кстати, по примерам из методички на простых запросах gorilla/mux выигрывает в скорости.
Вот результаты.
gorilla/mux:
https://mega.nz/file/5EojWKTT#cUbN5KHqAGGN1vDffN42X8Z41xuKMCtQ0Kgalga-c98

go-chi:
https://mega.nz/file/oVojDCRa#IVgixZfnfojZZ_Xy3aAKPtKgn7JFdek24TgCQfmuTLc

fasthttp:
https://mega.nz/file/VEhjVSyB#Q4ygmZn-_sAUABrD-GgGSrcZqmabl-x0101WhQW2Kcw