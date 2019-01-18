DROP TABLE IF EXISTS forum_user CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS thread CASCADE;
DROP TABLE IF EXISTS vote CASCADE;
DROP TABLE IF EXISTS user_in_forum CASCADE;

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS forum_user (
    email CITEXT UNIQUE, 
    about text,
    fullname text,
    nickname CITEXT COLLATE "ucs_basic" primary key 
);

CREATE INDEX IF NOT EXISTS forum_user_nickname_email_idx ON forum_user(email);

------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS forum (
    posts INTEGER,
    slug CITEXT primary key,
    threads INTEGER,
    title text,
    forum_user CITEXT NOT NULL,

    FOREIGN KEY (forum_user) REFERENCES forum_user(nickname) 
);

-------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS thread (
    author CITEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum CITEXT NOT NULL,
    id SERIAL primary key,
    isEdited BOOL DEFAULT false, 
    Msg text,
    slug CITEXT DEFAULT NULL,
    title text,
    votes INTEGER DEFAULT 0,
    
    FOREIGN KEY(author) REFERENCES forum_user(nickname),
    FOREIGN KEY(forum) REFERENCES forum(slug)
);

CREATE INDEX IF NOT EXISTS thread_slug_idx on thread(slug);


CREATE OR REPLACE FUNCTION incThreads() RETURNS TRIGGER AS
$$BEGIN
    UPDATE forum SET threads = threads + 1 WHERE slug = NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

DROP TRIGGER IF EXISTS incThreadsOnInsertThread on thread;

CREATE TRIGGER incThreadsOnInsertThread 
AFTER INSERT ON thread
FOR EACH ROW EXECUTE PROCEDURE incThreads();


-------------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS post (
    author CITEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE,
    forum CITEXT NOT NULL,
    id SERIAL primary key,
    isEdited BOOL DEFAULT false, 
    Msg text,
    parent INTEGER,
    thread INTEGER NOT NULL,
    mpath INTEGER ARRAY,

    FOREIGN KEY(author) REFERENCES forum_user(nickname),
    FOREIGN KEY(forum) REFERENCES forum(slug),
    FOREIGN KEY(thread) REFERENCES thread(id)
);

CREATE INDEX IF NOT EXISTS post_thread_idx ON post(thread, id);

CREATE OR REPLACE FUNCTION incPosts() RETURNS TRIGGER AS
$$BEGIN
    UPDATE forum SET posts = posts + 1 WHERE slug = NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

DROP TRIGGER IF EXISTS incPostsOnInsertPost on post;

CREATE TRIGGER incPostsOnInsertPost 
AFTER INSERT ON post
FOR EACH ROW EXECUTE PROCEDURE incPosts();

CREATE OR REPLACE FUNCTION getMpath(INTEGER) RETURNS INTEGER ARRAY
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


------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS vote (
    nickname CITEXT NOT NULL,
    voice SMALLINT,
    thread INTEGER NOT NULL,

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


---------------------------------------------------------------------------------------
CREATE TABLE user_in_forum (
    nickname CITEXT NOT NULL,
    forum CITEXT NOT NULL,

    FOREIGN KEY(nickname) REFERENCES forum_user(nickname),
    FOREIGN KEY(forum) REFERENCES forum(slug),
    PRIMARY KEY(nickname, forum)
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