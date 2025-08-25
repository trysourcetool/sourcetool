BEGIN;

DROP INDEX idx_session_type;
ALTER TABLE "session" DROP COLUMN "type";

DROP INDEX idx_session_type_name;
DROP INDEX idx_session_type_code;
DROP TABLE "session_type";

END;