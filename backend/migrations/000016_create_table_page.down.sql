BEGIN;

DROP TRIGGER IF EXISTS validate_page ON "page";
DROP TRIGGER IF EXISTS update_page_updated_at ON "page";
DROP TABLE "page";

END;
