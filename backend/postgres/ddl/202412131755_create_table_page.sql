-- +migrate Up

CREATE TABLE "page" (
  "id"               UUID          NOT NULL,
  "organization_id"  UUID          NOT NULL,
  "environment_id"   UUID          NOT NULL,
  "api_key_id"       UUID          NOT NULL,
  "name"             VARCHAR(255)  NOT NULL,
  "route"            VARCHAR(255)  NOT NULL,
  "path"             INTEGER[]     NOT NULL,
  "created_at"       TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"       TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("environment_id") REFERENCES "environment"("id") ON DELETE CASCADE,
  FOREIGN KEY ("api_key_id") REFERENCES "api_key"("id") ON DELETE CASCADE,
  UNIQUE("organization_id", "api_key_id", "route"),
  UNIQUE("organization_id", "api_key_id", "path"),
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_page_updated_at
    BEFORE UPDATE ON "page"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER validate_page
    BEFORE INSERT OR UPDATE ON "page"
    FOR EACH ROW
    EXECUTE FUNCTION validate_page();

-- +migrate Down

DROP TRIGGER IF EXISTS validate_page ON "page";
DROP TRIGGER IF EXISTS update_page_updated_at ON "page";
DROP TABLE "page";
