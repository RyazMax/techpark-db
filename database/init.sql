DROP TABLE IF EXISTS forum_user CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS thread CASCADE;
DROP TABLE IF EXISTS vote CASCADE;

CREATE TABLE IF NOT EXISTS forum_user (
    email text unique,
    about text,
    fullname text,
    nickname text primary key 
);

CREATE TABLE IF NOT EXISTS forum (
    posts BIGSERIAL,
    slug text primary key,
    threads INTEGER,
    title text,
    forum_user text,

    FOREIGN KEY (forum_user) REFERENCES forum_user(nickname) 
);

CREATE TABLE IF NOT EXISTS thread (
    author text,
    created DATE,
    forum text,
    id BIGSERIAL primary key,
    isEdited BOOL, 
    message text,
    slug text,
    title text,
    votes BIGSERIAL,
    
    FOREIGN KEY(author) REFERENCES forum_user(nickname),
    FOREIGN KEY(forum) REFERENCES forum(slug)
);

CREATE TABLE IF NOT EXISTS post (
    author text,
    created DATE,
    forum text,
    id BIGSERIAL primary key,
    isEdited BOOL, 
    message text,
    parent BIGSERIAL,
    thread BIGSERIAL,

    FOREIGN KEY(author) REFERENCES forum_user(nickname),
    FOREIGN KEY(forum) REFERENCES forum(slug),
    FOREIGN KEY(thread) REFERENCES thread(id)
);

CREATE TABLE IF NOT EXISTS vote (
    nickname text,
    voice SMALLINT,

    FOREIGN KEY(nickname) REFERENCES forum_user(nickname)
);

CREATE INDEX idx_forum_user_nickname_email ON forum_user(nickname, email);
