CREATE ROLE keyhub;
CREATE USER "keyhub-local" WITH PASSWORD 'keyhub';
GRANT keyhub TO "keyhub-local";