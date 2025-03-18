BEGIN;

CREATE TABLE "environment" (
  "id"              UUID          NOT NULL,
  "organization_id" UUID          NOT NULL,
  "name"            VARCHAR(255)  NOT NULL,
  "slug"            VARCHAR(255)  NOT NULL,
  "color"           VARCHAR(255)  NOT NULL,
  "created_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_environment_organization_slug ON "environment" ("organization_id", "slug");

CREATE TRIGGER update_environment_updated_at
    BEFORE UPDATE ON "environment"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

END;
