BEGIN TRANSACTION;

PRAGMA foreign_keys = ON;

CREATE TABLE securities (
    security_id INTEGER PRIMARY KEY,
    added_on TEXT DEFAULT CURRENT_DATE,
    ended_on TEXT
);

CREATE TABLE limits (
    limit_id INTEGER PRIMARY KEY,
    value INTEGER NOT NULL,
    name TEXT NOT NULL,
    day_lag INTEGER
);

CREATE TABLE alerts (
    alert_id INTEGER PRIMARY KEY,
    created_on TEXT DEFAULT CURRENT_DATE,
    security_id INTEGER,
    limit_id INTEGER,
    UNIQUE (created_on, security_id, limit_id),
    FOREIGN KEY (security_id) REFERENCES securities (security_id),
    FOREIGN KEY (limit_id) REFERENCES thresholds (limit_id)
);

COMMIT;