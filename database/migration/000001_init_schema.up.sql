CREATE TABLE "Users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "is_verified" boolean,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_login" timestamptz NOT NULL DEFAULT (now()),
  "profile_blurb" varchar,
  "banned" boolean
);

CREATE TABLE "DomainBans" (
  "user_id" bigint,
  "banned_from" varchar,
  "banned_by" bigint,
  "banned_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Comments" (
  "id" bigserial PRIMARY KEY,
  "path_id" bigint NOT NULL,
  "author" bigint NOT NULL,
  "content" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "parent" bigint DEFAULT 0,
  "hidden" boolean DEFAULT false,
  "removed" boolean DEFAULT false
);

CREATE TABLE "VoteRecords" (
  "comment_id" bigint,
  "category" varchar,
  "user_id" bigint,
  "value" int8,
  PRIMARY KEY ("comment_id", "category")
);

CREATE TABLE "Domains" (
  "id" varchar UNIQUE PRIMARY KEY NOT NULL
);

CREATE TABLE "Paths" (
  "id" bigserial PRIMARY KEY,
  "domain" varchar,
  "path" varchar
);

CREATE TABLE "DomainModeratorAssignments" (
  "id" bigserial PRIMARY KEY,
  "domain" varchar NOT NULL,
  "assigned_to" bigint NOT NULL,
  "assigned_at" timestamptz NOT NULL DEFAULT (now()),
  "assigned_by" bigint NOT NULL,
  "is_deactivation" boolean NOT NULL DEFAULT false
);

CREATE TABLE "GlobalModeratorAssignments" (
  "id" bigserial PRIMARY KEY,
  "assigned_to" bigint NOT NULL,
  "assigned_at" timestamptz NOT NULL DEFAULT (now()),
  "assigned_by" bigint NOT NULL,
  "is_deactivation" boolean DEFAULT false
);

CREATE TABLE "AdminAssignments" (
  "id" bigserial PRIMARY KEY,
  "assigned_to" bigint NOT NULL,
  "assigned_by" bigint NOT NULL,
  "assigned_at" timestamptz NOT NULL DEFAULT (now()),
  "is_deactivation" boolean DEFAULT false
);

CREATE TABLE "Logs" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "ip" varchar,
  "url" varchar,
  "at_time" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "CommentModerationActions" (
  "id" bigserial PRIMARY KEY,
  "taken_by" bigint NOT NULL,
  "comment_id" bigint NOT NULL,
  "reason" varchar,
  "taken_on" timestamptz,
  "set_hidden_to" boolean,
  "set_removed_to" boolean,
  "associated_report" bigint
);

CREATE TABLE "BanActions" (
  "id" bigserial PRIMARY KEY,
  "taken_by" bigint NOT NULL,
  "target_user" bigint NOT NULL,
  "reason" varchar,
  "taken_on" timestamptz,
  "domain" varchar,
  "set_banned_to" boolean
);

CREATE TABLE "CommentReports" (
  "id" bigserial PRIMARY KEY,
  "reporting_user" bigint NOT NULL,
  "comment" bigint NOT NULL,
  "reason" varchar,
  "action_taken" boolean,
  "time_created" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "VerificationCodes" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "verify_code" varchar,
  "created_on" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "PasswordResetCodes" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "verify_code" varchar,
  "created_on" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Feedbacks" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "type" varchar,
  "submitted_at" timestamptz NOT NULL DEFAULT (now()),
  "content" varchar,
  "hidden" boolean
);

CREATE INDEX ON "Users" ("username");

CREATE INDEX ON "Users" ("email");

CREATE INDEX ON "DomainBans" ("user_id");

CREATE INDEX ON "DomainBans" ("banned_by");

CREATE INDEX ON "DomainBans" ("banned_from");

CREATE INDEX ON "Comments" ("author");

CREATE INDEX ON "Comments" ("path_id");

CREATE UNIQUE INDEX ON "Paths" ("domain", "path");

CREATE INDEX ON "DomainModeratorAssignments" ("assigned_to");

CREATE INDEX ON "GlobalModeratorAssignments" ("assigned_to");

CREATE INDEX ON "AdminAssignments" ("assigned_to");

COMMENT ON COLUMN "Users"."password" IS 'Must be encrypted';

ALTER TABLE "DomainBans" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "Comments" ADD FOREIGN KEY ("path_id") REFERENCES "Paths" ("id");

ALTER TABLE "Comments" ADD FOREIGN KEY ("author") REFERENCES "Users" ("id");

ALTER TABLE "Comments" ADD FOREIGN KEY ("parent") REFERENCES "Comments" ("id");

ALTER TABLE "VoteRecords" ADD FOREIGN KEY ("comment_id") REFERENCES "Comments" ("id");

ALTER TABLE "VoteRecords" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "Paths" ADD FOREIGN KEY ("domain") REFERENCES "Domains" ("id");

ALTER TABLE "DomainModeratorAssignments" ADD FOREIGN KEY ("assigned_to") REFERENCES "Users" ("id");

ALTER TABLE "DomainModeratorAssignments" ADD FOREIGN KEY ("assigned_by") REFERENCES "Users" ("id");

ALTER TABLE "GlobalModeratorAssignments" ADD FOREIGN KEY ("assigned_to") REFERENCES "Users" ("id");

ALTER TABLE "GlobalModeratorAssignments" ADD FOREIGN KEY ("assigned_by") REFERENCES "Users" ("id");

ALTER TABLE "AdminAssignments" ADD FOREIGN KEY ("assigned_to") REFERENCES "Users" ("id");

ALTER TABLE "AdminAssignments" ADD FOREIGN KEY ("assigned_by") REFERENCES "Users" ("id");

ALTER TABLE "Logs" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "CommentModerationActions" ADD FOREIGN KEY ("taken_by") REFERENCES "Users" ("id");

ALTER TABLE "CommentModerationActions" ADD FOREIGN KEY ("comment_id") REFERENCES "Comments" ("id");

ALTER TABLE "CommentModerationActions" ADD FOREIGN KEY ("associated_report") REFERENCES "CommentReports" ("id");

ALTER TABLE "BanActions" ADD FOREIGN KEY ("taken_by") REFERENCES "Users" ("id");

ALTER TABLE "BanActions" ADD FOREIGN KEY ("target_user") REFERENCES "Users" ("id");

ALTER TABLE "BanActions" ADD FOREIGN KEY ("domain") REFERENCES "Domains" ("id");

ALTER TABLE "CommentReports" ADD FOREIGN KEY ("reporting_user") REFERENCES "Users" ("id");

ALTER TABLE "CommentReports" ADD FOREIGN KEY ("comment") REFERENCES "Comments" ("id");

ALTER TABLE "VerificationCodes" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "PasswordResetCodes" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "Feedbacks" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");
