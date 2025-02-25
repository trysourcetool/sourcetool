BEGIN;

CREATE TABLE "host_instance" (
  "id"              UUID          NOT NULL,
  "organization_id" UUID          NOT NULL,
  "api_key_id"      UUID          NOT NULL,
  "sdk_name"        VARCHAR(255)  NOT NULL,
  "sdk_version"     VARCHAR(255)  NOT NULL,
  "status"          INTEGER       NOT NULL,
  "created_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("api_key_id") REFERENCES "api_key"("id") ON DELETE CASCADE,
  FOREIGN KEY ("status") REFERENCES "host_instance_status"("code") ON DELETE RESTRICT,
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_host_instance_updated_at
    BEFORE UPDATE ON "host_instance"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_host_instance
    BEFORE INSERT OR UPDATE ON "host_instance"
    FOR EACH ROW
    EXECUTE FUNCTION validate_host_instance();

END;
