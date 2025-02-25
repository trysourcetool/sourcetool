BEGIN;

DROP TRIGGER IF EXISTS update_environment_updated_at ON "environment";
DROP TABLE "environment";

END;
