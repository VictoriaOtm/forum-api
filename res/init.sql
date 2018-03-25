CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS t_user;
DROP TABLE IF EXISTS t_forum;
DROP TABLE IF EXISTS t_thread;
DROP TABLE IF EXISTS t_forum_user;
DROP TABLE IF EXISTS t_posts;
DROP TABLE IF EXISTS t_vote;

CREATE TABLE t_user (
  nickname citext CONSTRAINT t_user_nickname_pk PRIMARY KEY,
  email    citext CONSTRAINT t_user_email_unique UNIQUE,
  fullname text,
  about    text
);

CREATE TABLE t_forum (
  slug    citext CONSTRAINT t_forum_slug_pk PRIMARY KEY,
  title   TEXT,
  f_user  TEXT,
  threads INT    DEFAULT 0,
  posts   BIGINT DEFAULT 0
);

CREATE TABLE t_thread (
  id      SERIAL,
  slug    citext,
  author  TEXT,
  created TIMESTAMPTZ DEFAULT current_timestamp,
  forum   TEXT,
  message TEXT,
  title   TEXT,
  votes   INT         DEFAULT 0
);

DROP FUNCTION IF EXISTS inc_thread;
CREATE FUNCTION inc_thread() RETURNS TRIGGER AS
$BODY$
BEGIN
  UPDATE t_forum
  SET
    threads = threads + 1
  WHERE slug = NEW.forum;
  RETURN NULL;
END;
$BODY$
LANGUAGE plpgsql;

CREATE TRIGGER t_thread_insert_trg AFTER INSERT ON t_thread
  FOR EACH ROW EXECUTE PROCEDURE inc_thread();

CREATE TABLE t_forum_user (
  slug     TEXT,
  nickname citext,
  email    TEXT,
  fullname TEXT,
  about    TEXT
);

CREATE UNIQUE INDEX t_forum_user_slug_nickname_uidx
  ON t_forum_user (slug, nickname);

CREATE TABLE t_posts (
  id BIGSERIAL,
  author TEXT,
  forum TEXT,
  thread INT,
  created TIMESTAMPTZ,
  is_edited BOOLEAN DEFAULT false,
  message TEXT,
  parent BIGINT,
  parents BIGINT[]
);

CREATE TABLE t_vote (
  thread_id INT,
  nickname TEXT,
  vote INT,
  CONSTRAINT t_vote_thread_id_nickname_idx UNIQUE (thread_id, nickname)
);

