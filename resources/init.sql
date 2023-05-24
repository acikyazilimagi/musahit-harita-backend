CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE "locations" (
    "id" serial primary key,
    "city_ysk_id" int,
    "city_name" varchar NOT NULL,
    "district_id" int,
    "district_ysk_id" int,
    "district_name" varchar NOT NULL,
    "county_ysk_id" int,
    "county_name" varchar,
    "neighbourhood_ysk_id" int,
    "neighbourhood_id" unique int,
    "neighbourhood_name" varchar,
    volunteer_data int,
    "threshold" int
);

ALTER TABLE
    "locations"
ADD
    CONSTRAINT "uq_locations_ysk_id" UNIQUE (
        "city_ysk_id",
        "district_ysk_id",
        "county_ysk_id",
        "neighbourhood_ysk_id"
    );

CREATE INDEX "idx_locations_neighbourhood_ysk_id" ON "locations" ("neighbourhood_ysk_id");

CREATE TABLE "buildings" (
    "id" serial primary key,
    "ysk_id" int unique,
    "name" varchar NOT NULL,
    "location_id" int NOT NULL REFERENCES "locations" ("id")
);

CREATE TYPE "source_type" AS ENUM (
    'ysk',
    'volunteer',
    'ballot_box',
    'volunteer_location_count'
);

CREATE TABLE "sources" (
    "id" serial primary key,
    "name" varchar unique NOT NULL,
    "type" "source_type" NOT NULL
);

CREATE TABLE "volunteers" (
    "id" serial primary key,
    volunteer_doc jsonb,
    "building_id" int REFERENCES "buildings" ("id"),
    "location_id" int REFERENCES "locations" ("id"),
    "confirmed" boolean,
    "source_id" int REFERENCES "sources" ("id")
);

CREATE TABLE "ballot_boxes" (
    "id" serial primary key,
    "box_no" int,
    "ysk_id" int unique,
    "building_id" int not null REFERENCES "buildings" ("id"),
    "location_id" int not null REFERENCES "locations" ("id"),
    "voter_count" int
);

CREATE TABLE "users" (
    "id" serial primary key,
    "login" varchar unique,
    "passwd" varchar NOT NULL,
    "source_id" int REFERENCES "sources" ("id")
);

CREATE TABLE "volunteer_counts" (
    "id" serial primary key,
    "count" int,
    "priority" int,
    "source_id" int REFERENCES "sources" ("id"),
    "building_id" int REFERENCES "buildings" ("id"),
    "location_id" int REFERENCES "locations" ("id"),
    "neighbourhood_id" int
);

CREATE INDEX "idx_volunteer_counts_neighbourhood_id" ON "volunteer_counts" ("neighbourhood_id");

CREATE INDEX "idx_volunteer_counts_location_id" ON "volunteer_counts" ("location_id");

CREATE INDEX idx_buildings_id ON "buildings" ("id");

CREATE INDEX idx_ballot_boxes_building_id ON "ballot_boxes" ("building_id");