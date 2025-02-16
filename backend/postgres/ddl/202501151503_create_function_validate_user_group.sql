-- +migrate Up

-- +migrate StatementBegin
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
-- +migrate StatementEnd

-- +migrate Down

DROP FUNCTION IF EXISTS validate_user_group();