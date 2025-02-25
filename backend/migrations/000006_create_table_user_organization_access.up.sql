BEGIN;

CREATE TABLE "user_organization_access" (
  "id"              UUID        NOT NULL,
  "user_id"         UUID        NOT NULL,
  "organization_id" UUID        NOT NULL,
  "role"            INTEGER     NOT NULL,
  "created_at"      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("role") REFERENCES "user_organization_role"("code") ON DELETE RESTRICT,
  UNIQUE("user_id", "organization_id"),
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_user_organization_access_updated_at
    BEFORE UPDATE ON "user_organization_access"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

END;
