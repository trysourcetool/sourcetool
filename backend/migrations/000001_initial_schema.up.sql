BEGIN;

-- Function to automatically update updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column() 
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- user table
CREATE TABLE "user" (
  "id"                 UUID          NOT NULL,
  "email"              VARCHAR(255)  NOT NULL,
  "first_name"         VARCHAR(255)  NOT NULL,
  "last_name"          VARCHAR(255)  NOT NULL,
  "refresh_token_hash" VARCHAR(255)  NOT NULL,
  "google_id"          VARCHAR(255)  NOT NULL,
  "created_at"         TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"         TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_user_email ON "user" ("email");
CREATE UNIQUE INDEX idx_user_refresh_token_hash ON "user" ("refresh_token_hash");
CREATE UNIQUE INDEX idx_user_google_id ON "user" ("google_id");

CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- organization table
CREATE TABLE "organization" (
  "id"         UUID          NOT NULL,
  "subdomain"  VARCHAR(255)  DEFAULT NULL,
  "created_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_organization_subdomain ON organization (subdomain) WHERE subdomain IS NOT NULL;

CREATE TRIGGER update_organization_updated_at
    BEFORE UPDATE ON "organization"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- user_organization_role table
CREATE TABLE "user_organization_role" (
  "code" INTEGER      NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  PRIMARY KEY ("code")
);

CREATE UNIQUE INDEX idx_user_organization_role_code ON "user_organization_role" ("code");
CREATE UNIQUE INDEX idx_user_organization_role_name ON "user_organization_role" ("name");

INSERT INTO "user_organization_role" ("code", "name") VALUES
  (0, 'unknown'),
  (1, 'admin'),
  (2, 'developer'),
  (3, 'member');

-- user_organization_access table
CREATE TABLE "user_organization_access" (
  "id"              UUID        NOT NULL,
  "user_id"         UUID        NOT NULL,
  "organization_id" UUID        NOT NULL,
  "role"            INTEGER     NOT NULL,
  "created_at"      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("role") REFERENCES "user_organization_role"("code") ON DELETE RESTRICT,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_user_organization_access_user_organization ON "user_organization_access" ("user_id", "organization_id");

CREATE TRIGGER update_user_organization_access_updated_at
    BEFORE UPDATE ON "user_organization_access"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- user_invitation table
CREATE TABLE "user_invitation" (
  "id"              UUID          NOT NULL,
  "organization_id" UUID          NOT NULL,
  "email"           VARCHAR(255)  NOT NULL,
  "role"            INTEGER       NOT NULL,
  "created_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("role") REFERENCES "user_organization_role"("code") ON DELETE RESTRICT,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_user_invitation_organization_email ON "user_invitation" ("organization_id", "email");

CREATE TRIGGER update_user_invitation_updated_at
    BEFORE UPDATE ON "user_invitation"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- environment table
CREATE TABLE "environment" (
  "id"              UUID          NOT NULL,
  "organization_id" UUID          NOT NULL,
  "name"            VARCHAR(255)  NOT NULL,
  "slug"            VARCHAR(255)  NOT NULL,
  "color"           VARCHAR(255)  NOT NULL,
  "created_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_environment_organization_slug ON "environment" ("organization_id", "slug");

CREATE TRIGGER update_environment_updated_at
    BEFORE UPDATE ON "environment"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- api_key table
CREATE OR REPLACE FUNCTION validate_api_key()
RETURNS TRIGGER AS $$
DECLARE
    environment_org_id UUID;
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM "user_organization_access" ua
        WHERE ua.user_id = NEW.user_id
        AND ua.organization_id = NEW.organization_id
    ) THEN
        RAISE EXCEPTION 'User % must belong to organization % to create an API key', NEW.user_id, NEW.organization_id;
    END IF;

    SELECT organization_id INTO environment_org_id
    FROM "environment"
    WHERE id = NEW.environment_id;

    IF environment_org_id != NEW.organization_id THEN
        RAISE EXCEPTION 'Environment % must belong to organization % to create an API key', NEW.environment_id, NEW.organization_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "api_key" (
  "id"              UUID          NOT NULL,
  "organization_id" UUID          NOT NULL,
  "environment_id"  UUID          NOT NULL,
  "user_id"         UUID          NOT NULL,
  "name"            VARCHAR(255)  NOT NULL,
  "key"             TEXT          NOT NULL,
  "created_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("environment_id") REFERENCES "environment"("id") ON DELETE RESTRICT,
  FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE RESTRICT,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_api_key_key ON "api_key" ("key");

CREATE TRIGGER update_api_key_updated_at
    BEFORE UPDATE ON "api_key"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_api_key
    BEFORE INSERT OR UPDATE ON "api_key"
    FOR EACH ROW
    EXECUTE FUNCTION validate_api_key();

-- host_instance table
CREATE OR REPLACE FUNCTION validate_host_instance()
RETURNS TRIGGER AS $$
DECLARE
    api_key_org_id UUID;
BEGIN
    SELECT organization_id INTO api_key_org_id
    FROM "api_key"
    WHERE id = NEW.api_key_id;

    IF api_key_org_id != NEW.organization_id THEN
        RAISE EXCEPTION 'API key % must belong to organization % to create a host instance', NEW.api_key_id, NEW.organization_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

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

CREATE TABLE "host_instance" (
  "id"              UUID          NOT NULL,
  "organization_id" UUID          NOT NULL,
  "api_key_id"      UUID          NOT NULL,
  "sdk_name"        VARCHAR(255)  NOT NULL,
  "sdk_version"     VARCHAR(255)  NOT NULL,
  "status"          INTEGER       NOT NULL,
  "created_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("api_key_id") REFERENCES "api_key"("id") ON DELETE CASCADE,
  FOREIGN KEY ("status") REFERENCES "host_instance_status"("code") ON DELETE RESTRICT,
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_host_instance_updated_at
    BEFORE UPDATE ON "host_instance"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_host_instance
    BEFORE INSERT OR UPDATE ON "host_instance"
    FOR EACH ROW
    EXECUTE FUNCTION validate_host_instance();

-- page table
CREATE OR REPLACE FUNCTION validate_page()
RETURNS TRIGGER AS $$
DECLARE
    environment_org_id UUID;
    api_key_org_id UUID;
BEGIN
    SELECT organization_id INTO environment_org_id
    FROM "environment"
    WHERE id = NEW.environment_id;

    IF environment_org_id != NEW.organization_id THEN
        RAISE EXCEPTION 'Environment % must belong to organization % to create a page', NEW.environment_id, NEW.organization_id;
    END IF;

    SELECT organization_id INTO api_key_org_id
    FROM "api_key"
    WHERE id = NEW.api_key_id;

    IF api_key_org_id != NEW.organization_id THEN
        RAISE EXCEPTION 'API key % must belong to organization % to create a page', NEW.api_key_id, NEW.organization_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "page" (
  "id"               UUID          NOT NULL,
  "organization_id"  UUID          NOT NULL,
  "environment_id"   UUID          NOT NULL,
  "api_key_id"       UUID          NOT NULL,
  "name"             VARCHAR(255)  NOT NULL,
  "route"            VARCHAR(255)  NOT NULL,
  "path"             INTEGER[]     NOT NULL,
  "created_at"       TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"       TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("environment_id") REFERENCES "environment"("id") ON DELETE CASCADE,
  FOREIGN KEY ("api_key_id") REFERENCES "api_key"("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_page_organization_api_key_route ON "page" ("organization_id", "api_key_id", "route");
CREATE UNIQUE INDEX idx_page_organization_api_key_path ON "page" ("organization_id", "api_key_id", "path");

CREATE TRIGGER update_page_updated_at
    BEFORE UPDATE ON "page"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_page
    BEFORE INSERT OR UPDATE ON "page"
    FOR EACH ROW
    EXECUTE FUNCTION validate_page();

-- session table
CREATE OR REPLACE FUNCTION validate_session()
RETURNS TRIGGER AS $$
DECLARE
    environment_org_id UUID;
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM "user_organization_access" ua
        WHERE ua.user_id = NEW.user_id
        AND ua.organization_id = NEW.organization_id
    ) THEN
        RAISE EXCEPTION 'User % must belong to organization % to create a session', NEW.user_id, NEW.organization_id;
    END IF;

    SELECT organization_id INTO environment_org_id
    FROM "environment"
    WHERE id = NEW.environment_id;

    IF environment_org_id != NEW.organization_id THEN
        RAISE EXCEPTION 'Environment % must belong to organization % to create a session', NEW.environment_id, NEW.organization_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "session" (
  "id"               UUID        NOT NULL,
  "organization_id"  UUID        NOT NULL,
  "user_id"          UUID        NOT NULL,
  "environment_id"   UUID        NOT NULL,
  "created_at"       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE,
  FOREIGN KEY ("environment_id") REFERENCES "environment"("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_session_updated_at
    BEFORE UPDATE ON "session"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER validate_session
    BEFORE INSERT OR UPDATE ON "session"
    FOR EACH ROW
    EXECUTE FUNCTION validate_session();

-- session_host_instance table
CREATE OR REPLACE FUNCTION validate_session_host_instance()
RETURNS TRIGGER AS $$
DECLARE
    api_key_environment_id UUID;
    session_environment_id UUID;
BEGIN
    SELECT environment_id INTO api_key_environment_id
    FROM "api_key" ak
    JOIN "host_instance" hi ON hi.api_key_id = ak.id
    WHERE hi.id = NEW.host_instance_id;

    SELECT environment_id INTO session_environment_id
    FROM "session"
    WHERE id = NEW.session_id;

    IF api_key_environment_id != session_environment_id THEN
        RAISE EXCEPTION 'Host instance (API key environment: %) and Session (environment: %) must belong to the same environment', api_key_environment_id, session_environment_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "session_host_instance" (
  "id"               UUID        NOT NULL,
  "session_id"       UUID        NOT NULL,
  "host_instance_id" UUID        NOT NULL,
  "created_at"       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("session_id") REFERENCES "session"("id") ON DELETE CASCADE,
  FOREIGN KEY ("host_instance_id") REFERENCES "host_instance"("id") ON DELETE CASCADE,
  UNIQUE("session_id", "host_instance_id"),
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_session_host_instance_updated_at
    BEFORE UPDATE ON "session_host_instance"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_session_host_instance
    BEFORE INSERT OR UPDATE ON "session_host_instance"
    FOR EACH ROW
    EXECUTE FUNCTION validate_session_host_instance();

-- group table
CREATE TABLE "group" (
  "id"               UUID          NOT NULL,
  "organization_id"  UUID          NOT NULL,
  "name"             VARCHAR(255)  NOT NULL,
  "slug"             VARCHAR(255)  NOT NULL,
  "created_at"       TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"       TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  UNIQUE("organization_id", "slug"),
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_group_updated_at
    BEFORE UPDATE ON "group"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- user_group table
CREATE OR REPLACE FUNCTION validate_user_group()
RETURNS TRIGGER AS $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM "user_organization_access" ua
    JOIN "group" g ON g.organization_id = ua.organization_id
    WHERE ua.user_id = NEW.user_id
    AND g.id = NEW.group_id
  ) THEN
    RAISE EXCEPTION 'User % and Group % must belong to the same organization', NEW.user_id, NEW.group_id;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "user_group" (
  "id"         UUID        NOT NULL,
  "user_id"    UUID        NOT NULL,
  "group_id"   UUID        NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE,
  FOREIGN KEY ("group_id") REFERENCES "group"("id") ON DELETE CASCADE,
  UNIQUE("user_id", "group_id"),
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_user_group_updated_at
    BEFORE UPDATE ON "user_group"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_user_group
    BEFORE INSERT OR UPDATE ON "user_group"
    FOR EACH ROW
    EXECUTE FUNCTION validate_user_group();

-- group_page table
CREATE OR REPLACE FUNCTION validate_group_page()
RETURNS TRIGGER AS $$
DECLARE
    group_org_id UUID;
    page_org_id UUID;
BEGIN
    SELECT organization_id INTO group_org_id
    FROM "group"
    WHERE id = NEW.group_id;

    SELECT organization_id INTO page_org_id
    FROM "page"
    WHERE id = NEW.page_id;

    IF group_org_id != page_org_id THEN
        RAISE EXCEPTION 'Group % and Page % must belong to the same organization', NEW.group_id, NEW.page_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "group_page" (
  "id"         UUID        NOT NULL,
  "group_id"   UUID        NOT NULL,
  "page_id"    UUID        NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("group_id") REFERENCES "group"("id") ON DELETE CASCADE,
  FOREIGN KEY ("page_id") REFERENCES "page"("id") ON DELETE CASCADE,
  UNIQUE("group_id", "page_id"),
  PRIMARY KEY ("id")
);

CREATE TRIGGER update_group_page_updated_at
    BEFORE UPDATE ON "group_page"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_group_page
    BEFORE INSERT OR UPDATE ON "group_page"
    FOR EACH ROW
    EXECUTE FUNCTION validate_group_page();

END;