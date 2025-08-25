BEGIN;

-- agent table

CREATE OR REPLACE FUNCTION validate_agent()
RETURNS TRIGGER AS $$
DECLARE
    environment_org_id UUID;
    api_key_org_id UUID;
BEGIN
    SELECT organization_id INTO environment_org_id
    FROM "environment"
    WHERE id = NEW.environment_id;

    IF environment_org_id != NEW.organization_id THEN
        RAISE EXCEPTION 'Environment % must belong to organization % to create an agent', NEW.environment_id, NEW.organization_id;
    END IF;

    SELECT organization_id INTO api_key_org_id
    FROM "api_key"
    WHERE id = NEW.api_key_id;

    IF api_key_org_id != NEW.organization_id THEN
        RAISE EXCEPTION 'API key % must belong to organization % to create an agent', NEW.api_key_id, NEW.organization_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "agent" (
  "id"              UUID          NOT NULL,
  "organization_id" UUID          NOT NULL,
  "environment_id"  UUID          NOT NULL,
  "api_key_id"      UUID          NOT NULL,
  "name"            VARCHAR(255)  NOT NULL,
  "description"     TEXT          NOT NULL DEFAULT '',
  "instructions"    TEXT          NOT NULL DEFAULT '',
  "model"           VARCHAR(255)  NOT NULL,
  "created_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE,
  FOREIGN KEY ("environment_id")  REFERENCES "environment"("id") ON DELETE CASCADE,
  FOREIGN KEY ("api_key_id")      REFERENCES "api_key"("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_agent_organization_api_key_name ON "agent" ("organization_id", "api_key_id", "name");

CREATE TRIGGER update_agent_updated_at
    BEFORE UPDATE ON "agent"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER validate_agent
    BEFORE INSERT OR UPDATE ON "agent"
    FOR EACH ROW
    EXECUTE FUNCTION validate_agent();

-- agent_tool table

CREATE TABLE "agent_tool" (
  "id"          UUID          NOT NULL,
  "agent_id"    UUID          NOT NULL,
  "name"        VARCHAR(255)  NOT NULL,
  "description" TEXT          NOT NULL DEFAULT '',
  "created_at"  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("agent_id") REFERENCES "agent"("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_agent_tool_agent_name ON "agent_tool" ("agent_id", "name");

CREATE TRIGGER update_agent_tool_updated_at
    BEFORE UPDATE ON "agent_tool"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

END;
