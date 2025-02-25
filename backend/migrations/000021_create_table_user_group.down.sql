BEGIN;

DROP TRIGGER IF EXISTS validate_user_group ON "user_group";
DROP TRIGGER IF EXISTS update_user_group_updated_at ON "user_group";
DROP TABLE "user_group";

END;
