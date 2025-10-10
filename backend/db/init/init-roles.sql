CREATE ROLE keyhub_role;
CREATE USER "keyhub-local" WITH PASSWORD 'keyhub';
GRANT keyhub_role TO "keyhub-local";