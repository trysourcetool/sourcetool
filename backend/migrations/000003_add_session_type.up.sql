BEGIN;

CREATE TABLE "session_type" (
  "code" INTEGER      NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  PRIMARY KEY ("code")
);

CREATE UNIQUE INDEX idx_session_type_code ON "session_type" ("code");
CREATE UNIQUE INDEX idx_session_type_name ON "session_type" ("name");

INSERT INTO "session_type" ("code", "name") VALUES
  (0, 'unknown'),
  (1, 'page'),
  (2, 'agent');

ALTER TABLE "session"
ADD COLUMN "type" INTEGER NOT NULL DEFAULT 1,
ADD FOREIGN KEY ("type") REFERENCES "session_type"("code") ON DELETE RESTRICT;

CREATE INDEX idx_session_type ON "session"("type");

END;