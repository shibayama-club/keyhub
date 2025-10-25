-- Create keyhub role with necessary permissions
CREATE ROLE keyhub WITH LOGIN PASSWORD 'keyhub';

-- Grant necessary permissions to keyhub on public schema
GRANT ALL ON SCHEMA public TO keyhub;
GRANT ALL ON ALL TABLES IN SCHEMA public TO keyhub;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO keyhub;
GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO keyhub;

-- For backward compatibility, create keyhub-local as an alias
CREATE USER "keyhub-local" WITH PASSWORD 'keyhub';
GRANT keyhub TO "keyhub-local";