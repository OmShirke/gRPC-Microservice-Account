-- Create the user account_om if it does not already exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'account_om') THEN
        CREATE USER account_om WITH PASSWORD 'POSTGRES_PASSWORD';
    END IF;
END
$$;

-- Grant all privileges on the database to the account_om user
GRANT ALL PRIVILEGES ON DATABASE postgres TO account_om;

-- Create the accounts table if it doesn't already exist
CREATE TABLE IF NOT EXISTS accounts (
  id CHAR(27) PRIMARY KEY,
  name VARCHAR(24) NOT NULL
);
