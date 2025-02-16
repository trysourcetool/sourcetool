-- +migrate Up

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION validate_session()
RETURNS TRIGGER AS $$
DECLARE
    page_org_id UUID;
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

    SELECT organization_id INTO page_org_id
    FROM "page"
    WHERE id = NEW.page_id;

    IF page_org_id != NEW.organization_id THEN
        RAISE EXCEPTION 'Page % must belong to organization % to create a session', NEW.page_id, NEW.organization_id;
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
-- +migrate StatementEnd

-- +migrate Down

DROP FUNCTION IF EXISTS validate_session();