DROP TABLE IF EXISTS forum_user;
CREATE TABLE IF NOT EXISTS forum_user (
    email text unique,
    about text,
    fullname text,
    nickname text primary key 
);