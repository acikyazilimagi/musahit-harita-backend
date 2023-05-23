CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE "locations"
(
    "id"                   serial primary key,
    "city_id"              int,
    "city_ysk_id"          int,
    "city_name"            varchar    NOT NULL,
    "district_id"          int,
    "district_ysk_id"      int,
    "district_name"        varchar    NOT NULL,
    "county_ysk_id"        int,
    "county_name"          varchar,
    "neighbourhood_id"     int unique NOT NULL,
    "neighbourhood_ysk_id" int,
    "neighbourhood_name"   varchar,
    volunteer_data         int,
    "threshold"            int
);

ALTER TABLE "locations"
    ADD CONSTRAINT "uq_locations_ysk_id" UNIQUE ("city_ysk_id", "district_ysk_id", "county_ysk_id",
                                                 "neighbourhood_ysk_id");

CREATE INDEX "idx_locations_neighbourhood_id" ON "locations" ("neighbourhood_id");

CREATE TABLE "buildings"
(
    "id"          serial primary key,
    "ysk_id"      int unique,
    "name"        varchar NOT NULL,
    "location_id" int     NOT NULL REFERENCES "locations" ("id")
);

CREATE TABLE "sources"
(
    "id"   serial primary key,
    "name" varchar unique NOT NULL
);

CREATE TABLE "volunteers"
(
    "id"          serial primary key,
    volunteer_doc jsonb,
    "building_id" int REFERENCES "buildings" ("id"),
    "location_id" int REFERENCES "locations" ("id"),
    "confirmed"   boolean,
    "source_id"   int REFERENCES "sources" ("id")
);

CREATE TABLE "ballot_boxes"
(
    "id"          serial primary key,
    "box_no"      int,
    "ysk_id"      int unique,
    "building_id" int not null REFERENCES "buildings" ("id"),
    "location_id" int not null REFERENCES "locations" ("id"),
    "voter_count" int
);

CREATE TABLE "users"
(
    "id"        serial primary key,
    "username"  varchar unique,
    "password"  varchar NOT NULL,
    "source_id" int REFERENCES "sources" ("id")
);

CREATE TABLE "volunteer_location_counts"
(
    "id"          serial primary key,
    "count"       int,
    "source_id"   int REFERENCES "sources" ("id"),
    "building_id" int REFERENCES "buildings" ("id"),
    "location_id" int unique REFERENCES "locations" ("id")
);