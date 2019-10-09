-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE "internals" (
    "id" bigserial,
    "name" varchar(255) NOT NULL,
    "key" varchar(255) NOT NULL,
    "role" varchar(255) NOT NULL,
    "secret" varchar(255) NOT NULL,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "messagings" (
    "id" bigserial,
    "partner" varchar(255) NOT NULL,
    "status" BOOLEAN,
    "send_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS "internals" CASCADE;
DROP TABLE IF EXISTS "messagings" CASCADE;