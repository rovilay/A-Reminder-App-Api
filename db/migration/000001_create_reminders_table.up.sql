CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE status AS ENUM ('incomplete', 'completed');

CREATE TABLE "reminders" (
    "id" UUID DEFAULT gen_random_uuid (),
    "message" VARCHAR NOT NULL,
    "time" TIMESTAMPTZ NOT NULL,
    "longitude" NUMERIC (8, 5) CHECK ("longitude" >= -180 AND "longitude" <= 180) DEFAULT NULL,
    "latitude" NUMERIC (7, 5) CHECK ("latitude" >= -90 AND "latitude" <= 90) DEFAULT NULL,
    "status" status NOT NULL DEFAULT 'incomplete',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW ()),
    "updated_at" TIMESTAMPTZ DEFAULT (NOW ()),
    PRIMARY KEY ("id")
);
