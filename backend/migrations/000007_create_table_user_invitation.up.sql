BEGIN;

CREATE TABLE "user_invitation" (
  "id"              UUID          NOT NULL,
  "organization_id" UUID          NOT NULL,
  "email"           VARCHAR(255)  NOT NULL UNIQUE,
  "role"            INTEGER       NOT NULL,
  "created_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("role") REFERENCES "user_organization_role"("code") ON DELETE RESTRICT,
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_user_invitation_updated_at
    BEFORE UPDATE ON "user_invitation"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

END;
