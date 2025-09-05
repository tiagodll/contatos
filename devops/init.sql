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
    [from]          TEXT NOT NULL,
    [to]            TEXT NOT NULL,
    [status]        TEXT NOT NULL DEFAULT 'requested',
    [message]       TEXT,
    [timestamp]     TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_friend_request UNIQUE ([from], [to])
);

CREATE INDEX IF NOT EXISTS idx_profiles_name ON profiles([name]);
