BEGIN;

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

END;
