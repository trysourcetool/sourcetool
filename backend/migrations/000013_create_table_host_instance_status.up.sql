BEGIN;

CREATE TABLE "host_instance_status" (
  "code" INTEGER      NOT NULL UNIQUE,
  "name" VARCHAR(255) NOT NULL UNIQUE,
  PRIMARY KEY ("code")
);

INSERT INTO "host_instance_status" ("code", "name") VALUES
  (0, 'unknown'),
  (1, 'online'),
  (2, 'unreachable'),
  (3, 'offline'),
  (4, 'shuttingDown');

END;
