BEGIN;

CREATE TABLE "user" (
  "id"         UUID          NOT NULL,
  "email"      VARCHAR(255)  NOT NULL,
  "first_name" VARCHAR(255)  NOT NULL,
  "last_name"  VARCHAR(255)  NOT NULL,
  "secret"     VARCHAR(255)  NOT NULL,
  "google_id"  VARCHAR(255)  NOT NULL,
  "created_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_user_email ON "user" ("email");
CREATE UNIQUE INDEX idx_user_secret ON "user" ("secret");

CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

END;
