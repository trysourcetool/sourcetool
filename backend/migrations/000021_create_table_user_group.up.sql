BEGIN;

CREATE TABLE "user_group" (
  "id"         UUID        NOT NULL,
  "user_id"    UUID        NOT NULL,
  "group_id"   UUID        NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE,
  FOREIGN KEY ("group_id") REFERENCES "group"("id") ON DELETE CASCADE,
  UNIQUE("user_id", "group_id"),
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_user_group_updated_at
    BEFORE UPDATE ON "user_group"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_user_group
    BEFORE INSERT OR UPDATE ON "user_group"
    FOR EACH ROW
    EXECUTE FUNCTION validate_user_group();

END;
