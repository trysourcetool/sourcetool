BEGIN;

CREATE TABLE "api_key" (
  "id"              UUID          NOT NULL,
  "organization_id" UUID          NOT NULL,
  "environment_id"  UUID          NOT NULL,
  "user_id"         UUID          NOT NULL,
  "name"            VARCHAR(255)  NOT NULL,
  "key"             TEXT          NOT NULL UNIQUE,
  "created_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("environment_id") REFERENCES "environment"("id") ON DELETE RESTRICT,
  FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE RESTRICT,
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_api_key_updated_at
    BEFORE UPDATE ON "api_key"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_api_key
    BEFORE INSERT OR UPDATE ON "api_key"
    FOR EACH ROW
    EXECUTE FUNCTION validate_api_key();

END;
