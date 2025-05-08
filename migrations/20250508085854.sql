-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "name" text NULL,
  "email" text NULL,
  "bio" text NULL,
  "member_number" text NULL,
  "avatar" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email")
);
-- Create "books" table
CREATE TABLE "public"."books" (
  "id" bigserial NOT NULL,
  "user_id" bigint NULL,
  "title" text NULL,
  "author" text NULL,
  "publisher" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_books" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "logins" table
CREATE TABLE "public"."logins" (
  "id" bigserial NOT NULL,
  "email" text NULL,
  "user_id" bigint NULL,
  "hash_password" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_logins_user_id" UNIQUE ("user_id"),
  CONSTRAINT "fk_users_login" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_logins_email" to table: "logins"
CREATE UNIQUE INDEX "idx_logins_email" ON "public"."logins" ("email");
