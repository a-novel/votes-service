CREATE TYPE vote AS ENUM ('up', 'down');

--bun:split

CREATE TABLE IF NOT EXISTS votes (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,

    vote vote NOT NULL,
    user_id uuid NOT NULL,
    target_id uuid NOT NULL,
    target TEXT NOT NULL,

    /* User is only allowed to vote once for a specific post. Targets ID is only unique within a certain target. */
    UNIQUE(user_id, target_id, target)
);

--bun:split

CREATE INDEX IF NOT EXISTS votes_summary_idx ON votes (target_id, target);
CREATE INDEX IF NOT EXISTS user_votes_idx ON votes (user_id, target);

--bun:split

CREATE VIEW votes_summary AS
    SELECT
        target_id,
        target,
        SUM(CASE WHEN vote = 'up' THEN 1 ELSE 0 END) AS up_votes,
        SUM(CASE WHEN vote = 'down' THEN 1 ELSE 0 END) AS down_votes
    FROM votes
    GROUP BY target_id, target;
