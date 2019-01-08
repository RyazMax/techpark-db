DROP TABLE IF EXISTS forum_user CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS thread CASCADE;
DROP TABLE IF EXISTS vote CASCADE;

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS forum_user (
    email text unique,
    about text,
    fullname text,
    nickname CITEXT COLLATE "ucs_basic" primary key 
);

CREATE TABLE IF NOT EXISTS forum (
    posts BIGINT,
    slug text primary key,
    threads INTEGER,
    title text,
    forum_user CITEXT,

    FOREIGN KEY (forum_user) REFERENCES forum_user(nickname) 
);

CREATE TABLE IF NOT EXISTS thread (
    author CITEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum text,
    id BIGSERIAL primary key,
    isEdited BOOL DEFAULT false, 
    Msg text,
    slug text DEFAULT NULL,
    title text,
    votes BIGINT DEFAULT 0,
    
    FOREIGN KEY(author) REFERENCES forum_user(nickname),
    FOREIGN KEY(forum) REFERENCES forum(slug)
);

/* Триггер на добавление thread
*/

CREATE OR REPLACE FUNCTION incThreads() RETURNS TRIGGER AS
$$BEGIN
    UPDATE forum SET threads = threads + 1 WHERE LOWER(slug) = LOWER(NEW.forum);
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

DROP TRIGGER IF EXISTS incThreadsOnInsertThread on thread;

CREATE TRIGGER incThreadsOnInsertThread 
AFTER INSERT ON thread
FOR EACH ROW EXECUTE PROCEDURE incThreads();

CREATE TABLE IF NOT EXISTS post (
    author CITEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum text,
    id BIGSERIAL primary key,
    isEdited BOOL DEFAULT false, 
    Msg text,
    parent BIGSERIAL,
    thread BIGSERIAL,
    mpath BIGINT ARRAY,

    FOREIGN KEY(author) REFERENCES forum_user(nickname),
    FOREIGN KEY(forum) REFERENCES forum(slug),
    FOREIGN KEY(thread) REFERENCES thread(id)
);
/*
Триггер на добавление поста
*/
CREATE OR REPLACE FUNCTION incPosts() RETURNS TRIGGER AS
$$BEGIN
    UPDATE forum SET posts = posts + 1 WHERE LOWER(slug) = LOWER(NEW.forum);
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

DROP TRIGGER IF EXISTS incPostsOnInsertPost on post;

CREATE TRIGGER incPostsOnInsertPost 
AFTER INSERT ON post
FOR EACH ROW EXECUTE PROCEDURE incPosts();



CREATE OR REPLACE FUNCTION getMpath(BIGINT) RETURNS BIGINT ARRAY
    AS 'SELECT mpath FROM post WHERE post.id = $1;'
    LANGUAGE SQL
    RETURNS NULL ON NULL INPUT;


DROP FUNCTION IF EXISTS addPath();

CREATE FUNCTION addPath() RETURNS TRIGGER AS
$$BEGIN
    NEW.mpath = array_append(getMpath(NEW.parent), NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

DROP TRIGGER IF EXISTS addPathOnInsertPost ON post;

CREATE TRIGGER addPathOnInsertPost
    BEFORE INSERT ON post
    FOR EACH ROW EXECUTE PROCEDURE addPath();

CREATE TABLE IF NOT EXISTS vote (
    nickname CITEXT,
    voice SMALLINT,
    thread BIGINT,

    FOREIGN KEY(nickname) REFERENCES forum_user(nickname),
    FOREIGN KEY(thread) REFERENCES thread(id)
);

CREATE OR REPLACE FUNCTION incVotes() RETURNS TRIGGER AS
$$BEGIN
    UPDATE thread SET votes = votes + NEW.voice WHERE id = NEW.thread;
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

DROP TRIGGER IF EXISTS incVotesOnInsertVote on vote;

CREATE OR REPLACE FUNCTION incVotesUpd() RETURNS TRIGGER AS
$$BEGIN
    UPDATE thread SET votes = votes + NEW.voice - OLD.voice WHERE id = NEW.thread;
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

CREATE TRIGGER incVotesOnInsertVote 
AFTER INSERT ON vote
FOR EACH ROW EXECUTE PROCEDURE incVotes();

DROP TRIGGER IF EXISTS incVotesOnUpdateVote on vote;

CREATE TRIGGER incVotesOnUpdateVote 
AFTER UPDATE ON vote
FOR EACH ROW EXECUTE PROCEDURE incVotesUpd();

DROP INDEX IF EXISTS forum_user_nickname_idx;
DROP INDEX IF EXISTS forum_user_nickname_email_idx;
DROP INDEX IF EXISTS vote_username_thread_idx;
DROP INDEX IF EXISTS thread_forum_created_idx;
DROP INDEX IF EXISTS thread_slug_idx;

DROP INDEX IF EXISTS posts_thread_idx;
DROP INDEX IF EXISTS posts_thread_created_idx;
DROP INDEX IF EXISTS post_mpath_idx;
DROP INDEX IF EXISTS post_mpath_desc_idx;

CREATE INDEX IF NOT EXISTS forum_user_nickname_idx ON forum_user(LOWER(nickname));
CREATE INDEX IF NOT EXISTS forum_user_nickname_email_idx ON forum_user(LOWER(nickname), LOWER(email));

CREATE INDEX IF NOT EXISTS thread_slug_idx on thread (LOWER(slug));
CREATE INDEX IF NOT EXISTS thread_forum_created_idx ON thread (LOWER(forum), created);
CREATE INDEX IF NOT EXISTS vote_username_thread_idx ON vote (LOWER(nickname), thread);

CREATE INDEX IF NOT EXISTS post_thread_idx ON post(thread);
CREATE INDEX IF NOT EXISTS posts_thread_created_idx ON post(thread, created);
CREATE INDEX IF NOT EXISTS post_mpath_idx ON post((mpath[1]), (mpath[2:]));
CREATE INDEX IF NOT EXISTS post_mpath_desc_id ON post((mpath[1]) desc, (mpath[2:]))