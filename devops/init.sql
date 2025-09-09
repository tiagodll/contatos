CREATE TABLE IF NOT EXISTS profiles
(
    [id]          TEXT NOT NULL COLLATE NOCASE,
    [name]        TEXT COLLATE NOCASE,
    [image]       TEXT COLLATE NOCASE,
    [data]        TEXT NOT NULL DEFAULT '{}',
    [timestamp]   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "PK_Profiles" PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS idx_profiles_name ON profiles([name]);

CREATE TABLE IF NOT EXISTS friend_request (
    [from]          TEXT NOT NULL,
    [to]            TEXT NOT NULL,
    [status]        TEXT NOT NULL DEFAULT 'requested',
    [message]       TEXT NOT NULL,
    [timestamp]     TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_friend_request UNIQUE ([from], [to])
);


CREATE TABLE IF NOT EXISTS friends (
    [user_id]          TEXT NOT NULL,
    [friend_id]            TEXT NOT NULL,
    [notes]         TEXT NOT NULL DEFAULT '',
    CONSTRAINT unique_friend_request UNIQUE ([user_id], [friend_id])
);
