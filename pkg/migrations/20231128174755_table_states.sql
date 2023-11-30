-- +goose Up
    CREATE TABLE IF NOT EXISTS states (
        id varchar(255),
        s_type int not null default 0,
        data text not null
    );

-- +goose StatementBegin
    CREATE OR REPLACE FUNCTION check_state()
        RETURNS TRIGGER AS $$
    DECLARE
        count int;
    BEGIN
        SELECT COUNT(*) INTO count FROM states WHERE id = NEW.id AND s_type = NEW.s_type;

        IF count > 0 THEN
           UPDATE states SET data = NEW.data WHERE id = NEW.id AND s_type = NEW.s_type;
           RETURN NULL;
        END IF;

        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;
-- +goose StatementEnd


    CREATE TRIGGER update_existing_state
        BEFORE INSERT ON states
        FOR EACH ROW
    EXECUTE FUNCTION check_state();

-- +goose Down