BEGIN;

DROP TRIGGER IF EXISTS validate_api_key ON "api_key";
DROP TRIGGER IF EXISTS update_api_key_updated_at ON "api_key";
DROP TABLE "api_key";

END;
