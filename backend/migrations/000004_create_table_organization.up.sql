BEGIN;

CREATE TABLE "organization" (
  "id"         UUID          NOT NULL,
  "subdomain"  VARCHAR(255)  DEFAULT NULL,
  "created_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_organization_subdomain ON organization (subdomain) WHERE subdomain IS NOT NULL;

CREATE TRIGGER update_organization_updated_at
    BEFORE UPDATE ON "organization"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

END;
