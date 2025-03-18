BEGIN;

CREATE TABLE "user_organization_role" (
  "code" INTEGER      NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  PRIMARY KEY ("code")
);

CREATE UNIQUE INDEX idx_user_organization_role_code ON "user_organization_role" ("code");
CREATE UNIQUE INDEX idx_user_organization_role_name ON "user_organization_role" ("name");

INSERT INTO "user_organization_role" ("code", "name") VALUES
  (0, 'unknown'),
  (1, 'admin'),
  (2, 'developer'),
  (3, 'member');

END;
