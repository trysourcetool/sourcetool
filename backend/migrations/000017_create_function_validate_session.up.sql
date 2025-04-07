BEGIN;

CREATE OR REPLACE FUNCTION validate_session()
RETURNS TRIGGER AS $$
DECLARE
    api_key_org_id UUID;
    host_instance_org_id UUID;
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM "user_organization_access" ua
        WHERE ua.user_id = NEW.user_id
        AND ua.organization_id = NEW.organization_id
    ) THEN
        RAISE EXCEPTION 'User % must belong to organization % to create a session', NEW.user_id, NEW.organization_id;
    END IF;

    SELECT organization_id INTO api_key_org_id
    FROM "api_key"
    WHERE id = NEW.api_key_id;

    IF api_key_org_id != NEW.organization_id THEN
        RAISE EXCEPTION 'API key % must belong to organization % to create a session', NEW.api_key_id, NEW.organization_id;
    END IF;

    SELECT organization_id INTO host_instance_org_id
    FROM "host_instance"
    WHERE id = NEW.host_instance_id;

    IF host_instance_org_id != NEW.organization_id THEN
        RAISE EXCEPTION 'Host instance % must belong to organization % to create a session', NEW.host_instance_id, NEW.organization_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

END;
