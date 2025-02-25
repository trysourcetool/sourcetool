BEGIN;

DROP TRIGGER IF EXISTS update_user_organization_access_updated_at ON "user_organization_access";
DROP TABLE "user_organization_access";

END;
