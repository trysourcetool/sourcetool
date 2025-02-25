BEGIN;

DROP TRIGGER IF EXISTS update_user_google_auth_request_updated_at ON "user_google_auth_request";
DROP TABLE "user_google_auth_request";

END;
