DROP TABLE IF EXISTS forum_user CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS thread CASCADE;
DROP TABLE IF EXISTS vote CASCADE;
DROP TABLE IF EXISTS user_in_forum CASCADE;

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS forum_user (
    email CITEXT unique,
    about text,
    fullname text,
    nickname CITEXT COLLATE "ucs_basic" primary key 
);

CREATE TABLE IF NOT EXISTS forum (
    posts BIGINT,
    slug CITEXT primary key,
    threads INTEGER,
    title text,
    forum_user CITEXT,

    FOREIGN KEY (forum_user) REFERENCES forum_user(nickname) 
);

CREATE TABLE IF NOT EXISTS thread (
    author CITEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum CITEXT,
    id BIGSERIAL primary key,
    isEdited BOOL DEFAULT false, 
    Msg text,
    slug CITEXT DEFAULT NULL,
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
    created TIMESTAMP WITH TIME ZONE,
    forum CITEXT,
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
    FOREIGN KEY(thread) REFERENCES thread(id),
    UNIQUE(nickname, thread)
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

CREATE TABLE user_in_forum (
    nickname CITEXT,
    forum CITEXT,

    FOREIGN KEY(nickname) REFERENCES forum_user(nickname),
    FOREIGN KEY(forum) REFERENCES forum(slug),
    UNIQUE(nickname,forum)
);


/* add user to forum*/
CREATE OR REPLACE FUNCTION addUserOnThreads() RETURNS TRIGGER AS
$$BEGIN
    INSERT INTO user_in_forum(nickname, forum) VALUES (NEW.author, NEW.forum) ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

DROP TRIGGER IF EXISTS addUserOnThreadsTRIG on thread;

CREATE TRIGGER addUserOnThreadsTRIG 
AFTER INSERT ON thread
FOR EACH ROW EXECUTE PROCEDURE addUserOnThreads();

CREATE OR REPLACE FUNCTION addUserOnPost() RETURNS TRIGGER AS
$$BEGIN
    INSERT INTO user_in_forum(nickname, forum) VALUES (NEW.author, NEW.forum) ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

DROP TRIGGER IF EXISTS addUserOnPostTRIG on post;

CREATE TRIGGER addUserOnPostTRIG
AFTER INSERT ON post
FOR EACH ROW EXECUTE PROCEDURE addUserOnPost();
 

DROP INDEX IF EXISTS user_in_forum_idx;
CREATE INDEX user_in_forum_idx on user_in_forum(forum, nickname);

DROP INDEX IF EXISTS forum_user_nickname_idx;
DROP INDEX IF EXISTS forum_user_nickname_email_idx;
DROP INDEX IF EXISTS vote_username_thread_idx;
DROP INDEX IF EXISTS thread_forum_created_idx;
DROP INDEX IF EXISTS thread_slug_idx;

DROP INDEX IF EXISTS posts_thread_idx;
DROP INDEX IF EXISTS posts_thread_created_idx;
DROP INDEX IF EXISTS post_mpath_idx;
DROP INDEX IF EXISTS post_mpath_desc_idx;

CREATE INDEX IF NOT EXISTS forum_user_nickname_idx ON forum_user(nickname);
CREATE INDEX IF NOT EXISTS forum_user_nickname_email_idx ON forum_user(nickname, email);

CREATE INDEX IF NOT EXISTS thread_slug_idx on thread (slug);
CREATE INDEX IF NOT EXISTS thread_forum_id on thread (forum);
CREATE INDEX IF NOT EXISTS thread_forum_created_idx ON thread (forum, created);
CREATE INDEX IF NOT EXISTS vote_username_thread_idx ON vote (nickname, thread);

--CREATE INDEX IF NOT EXISTS post_thread_idx ON post(thread, id);
--CREATE INDEX IF NOT EXISTS posts_thread_created_idx ON post(thread, created);
--CREATE INDEX IF NOT EXISTS post_mpath_idx ON post((mpath[1]))
--CREATE INDEX IF NOT EXISTS post_mpath_idx ON post((mpath[1]), (mpath[2:]));
--CREATE INDEX IF NOT EXISTS post_mpath_desc_id ON post((mpath[1]) desc, (mpath[2:]))