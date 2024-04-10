-- DO
-- $do$
-- BEGIN
--    IF EXISTS (SELECT FROM pg_database WHERE datname = 'account-db') THEN
--       RAISE NOTICE 'Database already exists';  -- optional
-- ELSE
--       PERFORM dblink_exec('dbname=' || current_database()  -- current db
--                         , 'CREATE DATABASE account-db');
--
-- END IF;
-- END
-- $do$;
--
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$
BEGIN
PERFORM dblink_exec('', 'CREATE DATABASE account-db');
EXCEPTION WHEN duplicate_database THEN RAISE NOTICE '%, skipping', SQLERRM USING ERRCODE = SQLSTATE;
END
$$;