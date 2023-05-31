---
date: 2023-05-30
refs:
  - blog/2022/gtfs.md
  - blog/2021/protobuf.md
---

# GTFS Realtime

![GTFS Realtime](gtfs-rt.png)

GTFS Realtime is a specification for providing realtime transit data.
This specification outlines the encoding of such data in the
[Protocol Buffers](../2021/protobuf.md) format, making it efficient.

Transit agencies publish and update a single file approximately once a minute,
which contains information about the arrival and departure of buses and trains
and their current location. For example, Google Maps, Apple Maps, and many
others use this information to power their transit directions.

[GTFS Realtime specification](https://gtfs.org/resources/gtfs-realtime/)  
[Mobility Database Catalogs](https://github.com/MobilityData/mobility-database-catalogs)

#transit
