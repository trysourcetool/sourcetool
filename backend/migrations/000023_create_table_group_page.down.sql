BEGIN;

DROP TRIGGER IF EXISTS validate_group_page ON "group_page";
DROP TRIGGER IF EXISTS update_group_page_updated_at ON "group_page";
DROP TABLE "group_page";

END;
