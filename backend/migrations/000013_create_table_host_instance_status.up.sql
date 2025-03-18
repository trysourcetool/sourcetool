BEGIN;

CREATE TABLE "host_instance_status" (
  "code" INTEGER      NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  PRIMARY KEY ("code")
);

CREATE UNIQUE INDEX idx_host_instance_status_code ON "host_instance_status" ("code");
CREATE UNIQUE INDEX idx_host_instance_status_name ON "host_instance_status" ("name");

INSERT INTO "host_instance_status" ("code", "name") VALUES
  (0, 'unknown'),
  (1, 'online'),
  (2, 'unreachable'),
  (3, 'offline'),
  (4, 'shuttingDown');

END;
