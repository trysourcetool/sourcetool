BEGIN;

DROP TRIGGER IF EXISTS validate_session ON "session";
DROP TRIGGER IF EXISTS update_session_updated_at ON "session";
DROP TABLE "session";

END;
