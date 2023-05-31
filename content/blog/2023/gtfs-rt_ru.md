---
date: 2023-05-30
refs:
  - blog/2022/gtfs_ru.md
  - blog/2021/protobuf_ru.md
---

# GTFS Realtime

![GTFS Realtime](gtfs-rt.png)

GTFS Realtime – это спецификация для передачи информации об общественном
транспорте в режиме реального времени. Эта спецификация описывает кодирование
таких данных в формате [Protocol Buffers](../2021/protobuf_ru.md).

Транспортные предприятия публикуют и примерно раз в минуту обновляют
единственный файл,в котором собрана информация о прибытии и отправлении
автобусов и поездов и их текущем местоположении. На основе этой информации
работают, например, маршруты в Google Maps, Apple Maps и многих других.

[GTFS Realtime specification](https://gtfs.org/resources/gtfs-realtime/)  
[Mobility Database Catalogs](https://github.com/MobilityData/mobility-database-catalogs)

#transit
