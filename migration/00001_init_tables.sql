-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE "clients" (
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
    "client_id" bigserial,
    "partner" varchar(125) NOT NULL,
    "status" varchar(16),
    "error_message" VARCHAR(255),
    "phone_number" varchar(17) NOT NULL,
    "message" varchar(153) NOT NULL,
    "raw_response" text NOT NULL,
    "send_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "notifications" (
    "id" bigserial,
    "client_id" bigserial,
    "title" varchar(255) NOT NULL,
    "message_body" text,
    "firebase_token" varchar(255),
    "topic" varchar(125),
    "response" varchar(255),
    "send_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "users" (
    "id" bigserial,
    "username" varchar(255) UNIQUE,
    "password" text NOT NULL,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS "clients" CASCADE;
DROP TABLE IF EXISTS "messagings" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "notifications" CASCADE;