CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS users_username_lower_key ON users (LOWER(username));
CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_key ON users (LOWER(email));

CREATE TABLE IF NOT EXISTS roles (
    name VARCHAR(20) PRIMARY KEY
);
INSERT INTO roles(name) VALUES ('admin'), ('staff'), ('user') ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS permissions (
    name VARCHAR(60) PRIMARY KEY
);
INSERT INTO permissions(name) VALUES
('item:list'), ('item:read:any'), ('item:create'), ('item:update:any'), ('item:delete'),
('claim:create'), ('claim:list'), ('claim:manage'), ('role:assign')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS role_permissions (
    role_name VARCHAR(20) NOT NULL REFERENCES roles(name) ON DELETE CASCADE,
    permission_name VARCHAR(60) NOT NULL REFERENCES permissions(name) ON DELETE CASCADE,
    PRIMARY KEY (role_name, permission_name)
);
INSERT INTO role_permissions(role_name, permission_name) VALUES
('admin','item:list'),('admin','item:read:any'),('admin','item:create'),('admin','item:update:any'),('admin','item:delete'),('admin','claim:create'),('admin','claim:list'),('admin','claim:manage'),('admin','role:assign'),
('staff','item:list'),('staff','item:read:any'),('staff','item:create'),('staff','item:update:any'),('staff','claim:create'),('staff','claim:list'),('staff','claim:manage'),
('user','item:list'),('user','item:create'),('user','item:delete'),('user','claim:create'),('user','claim:list')
ON CONFLICT DO NOTHING;

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_fkey;
UPDATE users SET role='user' WHERE role NOT IN (SELECT name FROM roles);
ALTER TABLE users ADD CONSTRAINT users_role_fkey FOREIGN KEY(role) REFERENCES roles(name);
CREATE INDEX IF NOT EXISTS users_role_idx ON users(role);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS refresh_tokens_user_idx ON refresh_tokens(user_id);

CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(120) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(60) NOT NULL,
    location VARCHAR(120) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'lost',
    reported_by BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT items_status_check CHECK(status IN ('lost','found','resolved'))
);
CREATE INDEX IF NOT EXISTS items_created_id_idx ON items(created_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS items_status_idx ON items(status);
CREATE INDEX IF NOT EXISTS items_category_idx ON items(category);

CREATE TABLE IF NOT EXISTS claims (
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    claimant_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    note TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT claims_status_check CHECK(status IN ('pending','approved','rejected')),
    UNIQUE(item_id, claimant_id)
);
CREATE INDEX IF NOT EXISTS claims_item_idx ON claims(item_id);
