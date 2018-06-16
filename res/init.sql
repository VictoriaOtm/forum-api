CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS t_user;
DROP TABLE IF EXISTS t_forum;
DROP TABLE IF EXISTS t_thread;
DROP TABLE IF EXISTS t_forum_user;
DROP TABLE IF EXISTS t_posts;
DROP TABLE IF EXISTS t_vote;

CREATE TABLE t_user (
  nickname citext PRIMARY KEY,
  email    citext UNIQUE,
  fullname text,
  about    text
);

CREATE INDEX t_user_idx_nickname
  on t_user (nickname, email, fullname, about);

CREATE TABLE t_forum (
  slug    citext PRIMARY KEY,
  title   TEXT,
  f_user  TEXT,
  threads INT    DEFAULT 0,
  posts   BIGINT DEFAULT 0
);

CREATE INDEX t_forum_idx_all
  on t_forum (slug, title, f_user, threads, posts);

CREATE TABLE t_thread (
  id      SERIAL PRIMARY KEY,
  slug    citext UNIQUE,
  author  TEXT,
  forum   citext,
  created TIMESTAMPTZ DEFAULT current_timestamp,
  message TEXT,
  title   TEXT,
  votes   INT         DEFAULT 0
);

CREATE INDEX t_thread_idx_forum_created
  ON t_thread (forum, created);

DROP FUNCTION IF EXISTS inc_thread();
CREATE FUNCTION inc_thread()
  RETURNS TRIGGER AS
$BODY$
BEGIN
  UPDATE t_forum SET threads = threads + 1 WHERE slug = NEW.forum;
  RETURN NULL;
END;
$BODY$
LANGUAGE plpgsql;

CREATE TRIGGER t_thread_insert_trg
  AFTER INSERT
  ON t_thread
  FOR EACH ROW EXECUTE PROCEDURE inc_thread();

CREATE TABLE t_forum_user (
  slug     citext,
  nickname citext,
  email    TEXT,
  fullname TEXT,
  about    TEXT
);

CREATE UNIQUE INDEX t_forum_user_slug_nickname_uidx
  ON t_forum_user (slug, nickname);

CREATE INDEX t_forum_user_slug_nickname_uidx_cov
  ON t_forum_user (slug, nickname, email, fullname, about);

CREATE TABLE t_posts (
  id          BIGSERIAL PRIMARY KEY,
  author      TEXT,
  forum       TEXT,
  thread      INT,
  created     TIMESTAMPTZ,
  is_edited   BOOLEAN DEFAULT false,
  message     TEXT,
  parent      BIGINT,
  main_parent BIGINT,
  parents     BIGINT []
);
--
-- CREATE INDEX t_posts_thread_id
--   ON t_posts (thread);

-- flat sort
CREATE INDEX t_posts_idx_thread_id
  ON t_posts (thread, id);

-- tree sort
CREATE INDEX t_posts_idx_thread_parents
  on t_posts (thread, parents);

-- tree sort
create index t_posts_idx_id_parents
  ON t_posts (id, parents);

create index t_posts_idx_id_parents2
  ON t_posts (id, main_parent);

-- parent tree sort
CREATE INDEX t_posts_idx_thread_parent_id
  on t_posts (thread, parent, id)
  WHERE parent = 0;

-- parent tree sort
CREATE INDEX t_posts_idx_thread_parent_main_parent_id
  on t_posts (thread, id, parent, main_parent)
  WHERE parent = 0;

-- parent tree sort
CREATE INDEX sssss
  ON t_posts (main_parent, parents);


CREATE TABLE t_vote (
  thread_id INT,
  nickname  TEXT NOT NULL,
  prev_vote SMALLINT DEFAULT 0,
  vote      SMALLINT,
  CONSTRAINT t_vote_thread_id_nickname_idx UNIQUE (thread_id, nickname)
);

