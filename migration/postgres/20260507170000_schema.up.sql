BEGIN;

CREATE TABLE tags
(
    id   TEXT PRIMARY KEY,
    name TEXT NOT NULL,

    CONSTRAINT tags_name_not_empty
        CHECK (length(trim(name)) > 0)
);

-- No public signup - accounts are created by seeding this table directly (see the seed
-- migration, which seeds exactly one - a temporary admin/12345 account, meant to be changed by
-- hand from /account, not enforced by the app). id is a plain SERIAL (not a slug/UUID like every
-- other table) because the app is single-admin for now and internal/domain.AdminUserID just
-- hardcodes "id 1" - no role column: with only ever one possible value it would carry no
-- information beyond the row existing at all. login is the sign-in identifier, email is an
-- optional profile field never used for auth.
CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    login         TEXT NOT NULL UNIQUE,
    email         TEXT UNIQUE,
    password_hash TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_login_not_empty
        CHECK (length(trim(login)) > 0),

    CONSTRAINT users_email_not_empty
        CHECK (email IS NULL OR length(trim(email)) > 0)
);

-- A logged-in browser session - see internal/domain.Session. token_hash is the SHA-256 hash of
-- the random token set as the session cookie; the raw token is never stored.
CREATE TABLE sessions
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);

CREATE TABLE sources
(
    id      TEXT PRIMARY KEY,
    user_id INTEGER REFERENCES users (id), -- NULL = global source, visible to everyone (every
                                            -- seeded source today) - no auth-scoped creation flow
                                            -- wires this up yet, but the FK is real.

    name        TEXT NOT NULL,
    feed_url    TEXT NOT NULL,
    description TEXT NOT NULL,
    tag_id      TEXT NOT NULL REFERENCES tags (id),

    badge       TEXT NOT NULL,
    badge_color TEXT NOT NULL,

    status         TEXT NOT NULL DEFAULT 'active',
    last_synced_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT sources_name_not_empty
        CHECK (length(trim(name)) > 0),

    CONSTRAINT sources_feed_url_not_empty
        CHECK (length(trim(feed_url)) > 0),

    CONSTRAINT sources_status_valid
        CHECK (status IN ('active', 'inactive'))
);

CREATE INDEX sources_user_id_idx ON sources (user_id);
CREATE INDEX sources_tag_id_idx ON sources (tag_id);

-- sources has no article_count column on purpose - with no RSS sync job to keep it fresh, a
-- stored count can only ever go stale. Callers COUNT(*) these rows instead (see
-- internal/adapter/postgres/source.List/Get).
CREATE TABLE articles
(
    id          TEXT PRIMARY KEY,
    source_id   TEXT NOT NULL REFERENCES sources (id) ON DELETE CASCADE,
    external_id TEXT NOT NULL, -- stable ID from the source's own feed - see adapter/feed.ToArticles

    title   TEXT NOT NULL,
    summary TEXT NOT NULL,
    url     TEXT NOT NULL,

    published_at TIMESTAMPTZ NOT NULL,
    unread       BOOLEAN     NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT articles_title_not_empty
        CHECK (length(trim(title)) > 0),

    CONSTRAINT articles_url_not_empty
        CHECK (length(trim(url)) > 0),

    -- Dedup key for RSS/Atom sync (see adapter/postgres/article.InsertBatch's ON CONFLICT) -
    -- scoped to source_id rather than global, since external_id is only ever unique within the
    -- feed it came from, not across sources.
    CONSTRAINT articles_source_id_external_id_key
        UNIQUE (source_id, external_id)
);

CREATE INDEX articles_source_id_idx ON articles (source_id);
CREATE INDEX articles_published_at_idx ON articles (published_at DESC);

COMMIT;
