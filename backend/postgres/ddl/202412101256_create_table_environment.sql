-- +migrate Up

CREATE TABLE "environment" (
  "id"              UUID          NOT NULL,
  "organization_id" UUID          NOT NULL,
  "name"            VARCHAR(255)  NOT NULL,
  "slug"            VARCHAR(255)  NOT NULL,
  "color"           VARCHAR(255)  NOT NULL,
  "created_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  UNIQUE("organization_id", "slug"),
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_environment_updated_at
    BEFORE UPDATE ON "environment"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down

DROP TRIGGER IF EXISTS update_environment_updated_at ON "environment";
DROP TABLE "environment";
