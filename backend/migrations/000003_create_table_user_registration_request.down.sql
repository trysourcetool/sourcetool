BEGIN;

DROP TRIGGER IF EXISTS update_user_registration_request_updated_at ON "user_registration_request";
DROP TABLE "user_registration_request";

END;
