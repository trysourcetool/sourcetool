BEGIN;

CREATE TABLE "user_google_auth_request" (
  "id"         UUID          NOT NULL,
  "email"      VARCHAR(255)  NOT NULL,
  "domain"     VARCHAR(255)  NOT NULL,
  "google_id"  VARCHAR(255)  NOT NULL,
  "invited"    BOOLEAN       NOT NULL,
  "expires_at" TIMESTAMPTZ   NOT NULL,
  "created_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_user_google_auth_request_updated_at
    BEFORE UPDATE ON "user_google_auth_request"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

END;
