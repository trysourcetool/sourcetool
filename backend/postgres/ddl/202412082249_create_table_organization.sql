-- +migrate Up

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

-- +migrate Down

DROP TRIGGER IF EXISTS update_organization_updated_at ON "organization";
DROP TABLE "organization";
