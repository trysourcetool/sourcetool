BEGIN;

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
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_page_organization_api_key_route ON "page" ("organization_id", "api_key_id", "route");
CREATE UNIQUE INDEX idx_page_organization_api_key_path ON "page" ("organization_id", "api_key_id", "path");

CREATE TRIGGER update_page_updated_at
    BEFORE UPDATE ON "page"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_page
    BEFORE INSERT OR UPDATE ON "page"
    FOR EACH ROW
    EXECUTE FUNCTION validate_page();

END;
