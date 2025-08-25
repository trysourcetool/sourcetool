BEGIN;

-- drop agent_tool
DROP TRIGGER IF EXISTS update_agent_tool_updated_at ON "agent_tool";
DROP INDEX IF EXISTS idx_agent_tool_agent_name;
DROP TABLE IF EXISTS "agent_tool";

-- drop agent
DROP TRIGGER IF EXISTS validate_agent ON "agent";
DROP TRIGGER IF EXISTS update_agent_updated_at ON "agent";
DROP INDEX IF EXISTS idx_agent_organization_api_key_name;
DROP TABLE IF EXISTS "agent";

-- drop function
DROP FUNCTION IF EXISTS validate_agent();

END;
