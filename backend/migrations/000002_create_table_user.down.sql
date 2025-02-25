BEGIN;

DROP TRIGGER IF EXISTS update_user_updated_at ON "user";
DROP TABLE "user";

END;
