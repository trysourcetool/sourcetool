BEGIN;

-- First drop all tables
DROP TABLE IF EXISTS "group_page";
DROP TABLE IF EXISTS "user_group";
DROP TABLE IF EXISTS "group";
DROP TABLE IF EXISTS "session";
DROP TABLE IF EXISTS "page";
DROP TABLE IF EXISTS "host_instance";
DROP TABLE IF EXISTS "host_instance_status";
DROP TABLE IF EXISTS "api_key";
DROP TABLE IF EXISTS "environment";
DROP TABLE IF EXISTS "user_invitation";
DROP TABLE IF EXISTS "user_organization_access";
DROP TABLE IF EXISTS "user_organization_role";
DROP TABLE IF EXISTS "organization";
DROP TABLE IF EXISTS "user";

-- Then drop all validation triggers
DROP TRIGGER IF EXISTS validate_group_page_trigger ON "group_page";
DROP TRIGGER IF EXISTS validate_user_group_trigger ON "user_group";
DROP TRIGGER IF EXISTS validate_session_trigger ON "session";
DROP TRIGGER IF EXISTS validate_page_trigger ON "page";
DROP TRIGGER IF EXISTS validate_host_instance_trigger ON "host_instance";
DROP TRIGGER IF EXISTS validate_api_key_trigger ON "api_key";

-- Drop all update_at triggers
DROP TRIGGER IF EXISTS update_group_page_updated_at ON "group_page";
DROP TRIGGER IF EXISTS update_user_group_updated_at ON "user_group";
DROP TRIGGER IF EXISTS update_group_updated_at ON "group";
DROP TRIGGER IF EXISTS update_session_updated_at ON "session";
DROP TRIGGER IF EXISTS update_page_updated_at ON "page";
DROP TRIGGER IF EXISTS update_host_instance_updated_at ON "host_instance";
DROP TRIGGER IF EXISTS update_host_instance_status_updated_at ON "host_instance_status";
DROP TRIGGER IF EXISTS update_api_key_updated_at ON "api_key";
DROP TRIGGER IF EXISTS update_environment_updated_at ON "environment";
DROP TRIGGER IF EXISTS update_user_invitation_updated_at ON "user_invitation";
DROP TRIGGER IF EXISTS update_user_organization_access_updated_at ON "user_organization_access";
DROP TRIGGER IF EXISTS update_user_organization_role_updated_at ON "user_organization_role";
DROP TRIGGER IF EXISTS update_organization_updated_at ON "organization";
DROP TRIGGER IF EXISTS update_user_updated_at ON "user";

-- Finally drop all functions
DROP FUNCTION IF EXISTS validate_group_page();
DROP FUNCTION IF EXISTS validate_user_group();
DROP FUNCTION IF EXISTS validate_session();
DROP FUNCTION IF EXISTS validate_page();
DROP FUNCTION IF EXISTS validate_host_instance();
DROP FUNCTION IF EXISTS validate_api_key();
DROP FUNCTION IF EXISTS update_updated_at_column();

END;