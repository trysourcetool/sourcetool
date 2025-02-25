BEGIN;

CREATE TABLE "group_page" (
  "id"         UUID        NOT NULL,
  "group_id"   UUID        NOT NULL,
  "page_id"    UUID        NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("group_id") REFERENCES "group"("id") ON DELETE CASCADE,
  FOREIGN KEY ("page_id") REFERENCES "page"("id") ON DELETE CASCADE,
  UNIQUE("group_id", "page_id"),
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_group_page_updated_at
    BEFORE UPDATE ON "group_page"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_group_page
    BEFORE INSERT OR UPDATE ON "group_page"
    FOR EACH ROW
    EXECUTE FUNCTION validate_group_page();

END;
