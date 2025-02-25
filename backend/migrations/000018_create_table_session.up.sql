BEGIN;

CREATE TABLE "session" (
  "id"               UUID        NOT NULL,
  "organization_id"  UUID        NOT NULL,
  "user_id"          UUID        NOT NULL,
  "page_id"          UUID        NOT NULL,
  "host_instance_id" UUID        NOT NULL,
  "created_at"       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE,
  FOREIGN KEY ("page_id") REFERENCES "page"("id") ON DELETE CASCADE,
  FOREIGN KEY ("host_instance_id") REFERENCES "host_instance"("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_session_updated_at
    BEFORE UPDATE ON "session"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER validate_session
    BEFORE INSERT OR UPDATE ON "session"
    FOR EACH ROW
    EXECUTE FUNCTION validate_session();

END;
