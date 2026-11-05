CREATE TABLE categories (
    name TEXT NOT NULL CONSTRAINT categories_pkey PRIMARY KEY CONSTRAINT categories_name_length CHECK (length(name) BETWEEN 1 AND 60),
    parent TEXT CONSTRAINT categories_parent_length CHECK (parent IS NULL OR length(parent) BETWEEN 1 AND 60),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON COLUMN categories.name IS 'Display name of the category, e.g. Groceries; also its identifier';
COMMENT ON COLUMN categories.parent IS 'Group the category belongs to, e.g. Food and eating out; NULL for a top-level category';
COMMENT ON COLUMN categories.created_at IS 'When the category was first seen';

CREATE TABLE movements (
    id BIGINT GENERATED ALWAYS AS IDENTITY CONSTRAINT movements_pkey PRIMARY KEY,
    ref TEXT CONSTRAINT movements_ref_key UNIQUE CONSTRAINT movements_ref_length CHECK (ref IS NULL OR length(ref) BETWEEN 1 AND 200),
    occurred_on DATE NOT NULL,
    merchant TEXT NOT NULL CONSTRAINT movements_merchant_length CHECK (length(merchant) BETWEEN 1 AND 120),
    merchant_key TEXT NOT NULL,
    amount BIGINT NOT NULL CONSTRAINT movements_amount_positive CHECK (amount > 0),
    category TEXT NOT NULL CONSTRAINT movements_category_fkey REFERENCES categories (name) ON UPDATE CASCADE,
    method TEXT NOT NULL,
    note TEXT NOT NULL DEFAULT '' CONSTRAINT movements_note_length CHECK (length(note) <= 200),
    source TEXT NOT NULL,
    edited BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON COLUMN movements.id IS 'Surrogate key';
COMMENT ON COLUMN movements.ref IS 'Stable reference built by the importer from the statement line; NULL for movements added by hand';
COMMENT ON COLUMN movements.occurred_on IS 'Day the money moved, in Bogotá time';
COMMENT ON COLUMN movements.merchant IS 'Merchant or person the money went to, as shown in the app';
COMMENT ON COLUMN movements.merchant_key IS 'The merchant lowercased with single spaces, computed by the app; what merchant rules match on';
COMMENT ON COLUMN movements.amount IS 'Amount in whole Colombian pesos, always positive';
COMMENT ON COLUMN movements.category IS 'Category the movement counts towards';
COMMENT ON COLUMN movements.method IS 'How it was paid: credit (card) or debit (bank transfer)';
COMMENT ON COLUMN movements.note IS 'Free-text note written by the user; empty when none';
COMMENT ON COLUMN movements.source IS 'Where it came from: statement (imported) or manual (added by hand)';
COMMENT ON COLUMN movements.edited IS 'True once the user changed an imported movement, so later imports leave it alone';
COMMENT ON COLUMN movements.deleted_at IS 'When the user deleted it; NULL while it counts. Kept so imports do not bring it back';
COMMENT ON COLUMN movements.created_at IS 'When the row was created';
COMMENT ON COLUMN movements.updated_at IS 'When the row last changed';

CREATE INDEX movements_occurred_on_idx ON movements (occurred_on) WHERE deleted_at IS NULL;
CREATE INDEX movements_merchant_key_idx ON movements (merchant_key);
CREATE INDEX movements_category_idx ON movements (category);

CREATE TABLE merchant_rules (
    merchant_key TEXT NOT NULL CONSTRAINT merchant_rules_pkey PRIMARY KEY,
    category TEXT NOT NULL CONSTRAINT merchant_rules_category_fkey REFERENCES categories (name) ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON COLUMN merchant_rules.merchant_key IS 'Merchant name, lowercased with single spaces, that the rule matches';
COMMENT ON COLUMN merchant_rules.category IS 'Category every movement from this merchant gets';
COMMENT ON COLUMN merchant_rules.updated_at IS 'When the rule was last set';

CREATE INDEX merchant_rules_category_idx ON merchant_rules (category);

CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY CONSTRAINT users_pkey PRIMARY KEY,
    username TEXT NOT NULL CONSTRAINT users_username_key UNIQUE CONSTRAINT users_username_length CHECK (length(username) BETWEEN 1 AND 60),
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON COLUMN users.id IS 'Surrogate key';
COMMENT ON COLUMN users.username IS 'Login name';
COMMENT ON COLUMN users.password_hash IS 'bcrypt hash of the password';
COMMENT ON COLUMN users.created_at IS 'When the user was created';

CREATE UNIQUE INDEX users_single_key ON users ((true));

CREATE TABLE sessions (
    token_hash TEXT NOT NULL CONSTRAINT sessions_pkey PRIMARY KEY,
    user_id BIGINT NOT NULL CONSTRAINT sessions_user_id_fkey REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL
);

COMMENT ON COLUMN sessions.token_hash IS 'SHA-256 of the session token, hex; the token itself is only in the browser cookie';
COMMENT ON COLUMN sessions.user_id IS 'User the session belongs to';
COMMENT ON COLUMN sessions.expires_at IS 'When the session stops working';

CREATE INDEX sessions_user_id_idx ON sessions (user_id);
