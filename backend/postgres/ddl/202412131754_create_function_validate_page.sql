-- +migrate Up

-- +migrate StatementBegin
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
-- +migrate StatementEnd

-- +migrate Down

DROP FUNCTION IF EXISTS validate_page();