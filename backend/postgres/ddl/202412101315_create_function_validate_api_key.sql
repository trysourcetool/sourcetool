-- +migrate Up

-- +migrate StatementBegin
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
-- +migrate StatementEnd

-- +migrate Down

DROP FUNCTION IF EXISTS validate_api_key();