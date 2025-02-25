BEGIN;

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

END;
