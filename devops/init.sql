CREATE TABLE IF NOT EXISTS profiles
(
    [id]          TEXT NOT NULL COLLATE NOCASE,
    [name]        TEXT COLLATE NOCASE,
    [image]       TEXT COLLATE NOCASE,
    [data]        TEXT NOT NULL DEFAULT '{}',
    [timestamp]   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "PK_Profiles" PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS friend_request (
    [id]            TEXT PRIMARY KEY,
    [user_id]       TEXT NOT NULL,
    [friend_id]     TEXT NOT NULL,
    [status]        TEXT NOT NULL DEFAULT 'requested',
    [message]       TEXT,
    [created_at]    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_friend_request UNIQUE ([user_id], [friend_id])
);

CREATE TABLE IF NOT EXISTS friends (
    [id]            TEXT PRIMARY KEY,
    [user_id]       TEXT NOT NULL,
    [friend_id]     TEXT NOT NULL,
    [notes]         TEXT,
    [updated_at]    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_friend UNIQUE ([user_id], [friend_id])
);

CREATE INDEX IF NOT EXISTS idx_friend_request_to ON friend_request([user_id]);
CREATE INDEX IF NOT EXISTS idx_friend_request_from ON friend_request([friend_id]);
CREATE INDEX IF NOT EXISTS idx_profiles_name ON profiles([first_name], [last_name]);
