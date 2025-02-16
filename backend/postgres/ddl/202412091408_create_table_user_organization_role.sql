-- +migrate Up

CREATE TABLE "user_organization_role" (
  "code" INTEGER      NOT NULL UNIQUE,
  "name" VARCHAR(255) NOT NULL UNIQUE,
  PRIMARY KEY ("code")
);

INSERT INTO "user_organization_role" ("code", "name") VALUES
  (0, 'unknown'),
  (1, 'admin'),
  (2, 'developer'),
  (3, 'member');

-- +migrate Down

DROP TABLE "user_organization_role";
