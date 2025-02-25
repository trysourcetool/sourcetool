BEGIN;

DROP TRIGGER IF EXISTS update_group_updated_at ON "group";
DROP TABLE "group";

END;
