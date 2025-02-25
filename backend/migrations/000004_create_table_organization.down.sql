BEGIN;

DROP TRIGGER IF EXISTS update_organization_updated_at ON "organization";
DROP TABLE "organization";

END;
