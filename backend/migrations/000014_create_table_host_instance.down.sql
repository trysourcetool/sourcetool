BEGIN;

DROP TRIGGER IF EXISTS validate_host_instance ON "host_instance";
DROP TRIGGER IF EXISTS update_host_instance_updated_at ON "host_instance";
DROP TABLE "host_instance";

END;
