BEGIN;

CREATE TABLE "organization" (
  "id"         UUID          NOT NULL,
  "subdomain"  VARCHAR(255)  NOT NULL UNIQUE,
  "created_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_organization_updated_at
    BEFORE UPDATE ON "organization"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

END;
