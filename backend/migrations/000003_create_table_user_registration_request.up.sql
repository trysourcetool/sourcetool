BEGIN;

CREATE TABLE "user_registration_request" (
  "id"         UUID          NOT NULL,
  "email"      VARCHAR(255)  NOT NULL UNIQUE,
  "created_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);


CREATE TRIGGER update_user_registration_request_updated_at
    BEFORE UPDATE ON "user_registration_request"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

END;
