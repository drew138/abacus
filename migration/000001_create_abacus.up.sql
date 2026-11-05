CREATE TABLE categories (
    -- Display name of the category, e.g. "Groceries"; also its identifier
    name TEXT NOT NULL CONSTRAINT categories_name_key PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 60),
    -- Name of the group the category belongs to, e.g. "Food and eating out"; NULL for a top-level category
    parent TEXT CHECK (parent IS NULL OR length(parent) BETWEEN 1 AND 60),
    -- When the category was first seen, RFC 3339 UTC
    created_at TEXT NOT NULL
);

CREATE TABLE movements (
    -- Surrogate key
    id INTEGER PRIMARY KEY,
    -- Stable reference built by the importer from the statement line; NULL for movements added by hand
    ref TEXT CONSTRAINT movements_ref_key UNIQUE CHECK (ref IS NULL OR length(ref) BETWEEN 1 AND 200),
    -- Day the money moved, YYYY-MM-DD in Bogotá time
    occurred_on TEXT NOT NULL CHECK (occurred_on GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]'),
    -- Merchant or person the money went to, as shown in the app
    merchant TEXT NOT NULL CHECK (length(merchant) BETWEEN 1 AND 120),
    -- The merchant lowercased with single spaces, computed by the app; what merchant rules match on
    merchant_key TEXT NOT NULL,
    -- Amount in whole Colombian pesos, always positive
    amount INTEGER NOT NULL CHECK (amount > 0),
    -- Category the movement counts towards
    category TEXT NOT NULL CONSTRAINT movements_category_fkey REFERENCES categories (name) ON UPDATE CASCADE,
    -- How it was paid: credit (card) or debit (account)
    method TEXT NOT NULL,
    -- Free-text note written by the user; empty when none
    note TEXT NOT NULL DEFAULT '' CHECK (length(note) <= 200),
    -- Where it came from: statement (imported) or manual (added by hand)
    source TEXT NOT NULL,
    -- 1 once the user changed an imported movement, so later imports leave it alone
    edited INTEGER NOT NULL DEFAULT 0 CHECK (edited IN (0, 1)),
    -- When the user deleted it, RFC 3339 UTC; NULL while it counts. Kept so imports do not bring it back
    deleted_at TEXT,
    -- When the row was created, RFC 3339 UTC
    created_at TEXT NOT NULL,
    -- When the row last changed, RFC 3339 UTC
    updated_at TEXT NOT NULL
);

CREATE INDEX movements_occurred_on_idx ON movements (occurred_on) WHERE deleted_at IS NULL;
CREATE INDEX movements_merchant_key_idx ON movements (merchant_key);
CREATE INDEX movements_category_idx ON movements (category);

CREATE TABLE merchant_rules (
    -- Merchant name, lowercased with single spaces, that the rule matches
    merchant_key TEXT NOT NULL CONSTRAINT merchant_rules_merchant_key_key PRIMARY KEY,
    -- Category every movement from this merchant gets
    category TEXT NOT NULL CONSTRAINT merchant_rules_category_fkey REFERENCES categories (name) ON UPDATE CASCADE,
    -- When the rule was last set, RFC 3339 UTC
    updated_at TEXT NOT NULL
);

CREATE TABLE users (
    -- Surrogate key
    id INTEGER PRIMARY KEY,
    -- Login name
    username TEXT NOT NULL CONSTRAINT users_username_key UNIQUE CHECK (length(username) BETWEEN 1 AND 60),
    -- bcrypt hash of the password
    password_hash TEXT NOT NULL,
    -- When the user was created, RFC 3339 UTC
    created_at TEXT NOT NULL
);

CREATE TABLE sessions (
    -- SHA-256 of the session token, hex; the token itself is only in the browser's cookie
    token_hash TEXT NOT NULL CONSTRAINT sessions_token_hash_key PRIMARY KEY,
    -- User the session belongs to
    user_id INTEGER NOT NULL CONSTRAINT sessions_user_id_fkey REFERENCES users (id) ON DELETE CASCADE,
    -- When the session stops working, RFC 3339 UTC
    expires_at TEXT NOT NULL
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);
