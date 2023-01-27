CREATE TABLE "Users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_login" timestamptz NOT NULL DEFAULT (now()),
  "profile_blurb" varchar,
  "banned" boolean
);

CREATE TABLE "DomainBans" (
  "user" bigint,
  "banned_from" varchar,
  "banned_by" bigint,
  "banned_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Comments" (
  "id" bigserial PRIMARY KEY,
  "pathid" bigint NOT NULL,
  "author" bigint NOT NULL,
  "content" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "parent" bigint,
  "hidden" boolean DEFAULT false,
  "removed" boolean DEFAULT false
);

CREATE TABLE "VoteRecord" (
  "commentId" bigint,
  "category" varchar,
  "userId" bigint,
  "value" int8,
  PRIMARY KEY ("commentId", "category")
);

CREATE TABLE "Domains" (
  "id" varchar UNIQUE PRIMARY KEY NOT NULL
);

CREATE TABLE "Paths" (
  "id" bigserial PRIMARY KEY,
  "domain" varchar,
  "path" varchar
);

CREATE TABLE "DomainModerators" (
  "id" bigint,
  "domain" varchar NOT NULL,
  "user" bigint NOT NULL,
  "granted_at" timestamptz NOT NULL DEFAULT (now()),
  "granted_by" bigint NOT NULL
);

CREATE TABLE "GlobalModerators" (
  "id" bigint,
  "user" bigint PRIMARY KEY,
  "granted_at" timestamptz NOT NULL DEFAULT (now()),
  "granted_by" bigint NOT NULL
);

CREATE TABLE "Admins" (
  "id" bigserial PRIMARY KEY,
  "user" bigint
);

CREATE TABLE "Logs" (
  "id" bigserial PRIMARY KEY,
  "user" bigint,
  "ip" varchar,
  "url" varchar
);

CREATE TABLE "CommentModerationActions" (
  "id" bigserial PRIMARY KEY,
  "taken_by" bigint,
  "commentId" bigint,
  "reason" varchar,
  "taken_on" timestamptz,
  "set_hidden_to" boolean,
  "set_removed_to" boolean,
  "associated_report" bigint
);

CREATE TABLE "BanActions" (
  "id" bigserial PRIMARY KEY,
  "taken_by" bigint,
  "target_user" bigint,
  "reason" varchar,
  "taken_on" timestamptz,
  "domain" varchar,
  "set_banned_to" boolean
);

CREATE TABLE "Reports" (
  "id" bigserial PRIMARY KEY,
  "reporting_user" bigint,
  "comment" bigint,
  "reason" varchar,
  "action_taken" boolean
);

CREATE INDEX ON "Users" ("username");

CREATE INDEX ON "Users" ("email");

CREATE INDEX ON "DomainBans" ("user");

CREATE INDEX ON "DomainBans" ("banned_by");

CREATE INDEX ON "DomainBans" ("banned_from");

CREATE INDEX ON "Comments" ("author");

CREATE INDEX ON "Comments" ("pathid");

CREATE UNIQUE INDEX ON "Paths" ("domain", "path");

COMMENT ON COLUMN "Users"."password" IS 'Must be encrypted';

ALTER TABLE "DomainBans" ADD FOREIGN KEY ("user") REFERENCES "Users" ("id");

ALTER TABLE "Comments" ADD FOREIGN KEY ("pathid") REFERENCES "Paths" ("id");

ALTER TABLE "Comments" ADD FOREIGN KEY ("author") REFERENCES "Users" ("id");

ALTER TABLE "VoteRecord" ADD FOREIGN KEY ("commentId") REFERENCES "Comments" ("id");

ALTER TABLE "VoteRecord" ADD FOREIGN KEY ("userId") REFERENCES "Users" ("id");

ALTER TABLE "Paths" ADD FOREIGN KEY ("domain") REFERENCES "Domains" ("id");

ALTER TABLE "DomainModerators" ADD FOREIGN KEY ("domain") REFERENCES "Domains" ("id");

ALTER TABLE "DomainModerators" ADD FOREIGN KEY ("user") REFERENCES "Users" ("id");

ALTER TABLE "DomainModerators" ADD FOREIGN KEY ("granted_by") REFERENCES "Users" ("id");

ALTER TABLE "GlobalModerators" ADD FOREIGN KEY ("user") REFERENCES "Users" ("id");

ALTER TABLE "GlobalModerators" ADD FOREIGN KEY ("granted_by") REFERENCES "Users" ("id");

ALTER TABLE "Admins" ADD FOREIGN KEY ("user") REFERENCES "Users" ("id");

ALTER TABLE "Logs" ADD FOREIGN KEY ("user") REFERENCES "Users" ("id");

ALTER TABLE "CommentModerationActions" ADD FOREIGN KEY ("taken_by") REFERENCES "Users" ("id");

ALTER TABLE "CommentModerationActions" ADD FOREIGN KEY ("commentId") REFERENCES "Comments" ("id");

ALTER TABLE "CommentModerationActions" ADD FOREIGN KEY ("associated_report") REFERENCES "Reports" ("id");

ALTER TABLE "BanActions" ADD FOREIGN KEY ("taken_by") REFERENCES "Users" ("id");

ALTER TABLE "BanActions" ADD FOREIGN KEY ("target_user") REFERENCES "Users" ("id");

ALTER TABLE "BanActions" ADD FOREIGN KEY ("domain") REFERENCES "Domains" ("id");

ALTER TABLE "Reports" ADD FOREIGN KEY ("reporting_user") REFERENCES "Users" ("id");

ALTER TABLE "Reports" ADD FOREIGN KEY ("comment") REFERENCES "Comments" ("id");