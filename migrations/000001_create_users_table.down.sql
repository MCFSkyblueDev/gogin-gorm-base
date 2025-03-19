-- Drop tables in reverse order to avoid constraint errors
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS logins;
DROP TABLE IF EXISTS users;

DROP TRIGGER IF EXISTS users_update_timestamp ON users;
DROP TRIGGER IF EXISTS logins_update_timestamp ON logins;
DROP TRIGGER IF EXISTS books_update_timestamp ON books;
DROP FUNCTION IF EXISTS update_timestamp;


