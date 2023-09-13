DROP VIEW IF EXISTS votes_summary;

--bun:split

DROP INDEX IF EXISTS votes_summary_idx;
DROP INDEX IF EXISTS user_votes_idx;

--bun:split

DROP TABLE IF EXISTS votes;

--bun:split

DROP TYPE IF EXISTS vote;
