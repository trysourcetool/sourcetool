BEGIN;

DROP TRIGGER IF EXISTS update_user_invitation_updated_at ON "user_invitation";
DROP TABLE "user_invitation";

END;
