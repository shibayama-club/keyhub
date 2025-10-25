-- Create keyhub role with necessary permissions
CREATE ROLE keyhub WITH LOGIN PASSWORD 'keyhub';

GRANT USAGE, CREATE ON SCHEMA public TO keyhub;

-- For backward compatibility, create keyhub-local as an alias
CREATE USER "keyhub-local" WITH PASSWORD 'keyhub';
GRANT keyhub TO "keyhub-local";