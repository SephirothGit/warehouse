-- Warehouse Tracker — database schema

-- Users & roles (RBAC)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO roles (name) VALUES ('admin'), ('manager'), ('warehouse_worker');

CREATE TABLE user_roles (
    user_id INTEGER NOT NULL REFERENCES users(id),
    role_id INTEGER NOT NULL REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);

-- Auth: refresh tokens (rotation + reuse detection)
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    token_hash VARCHAR(64) NOT NULL,
    issued_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE
);

-- Warehouse hierarchy
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255)
);

CREATE TABLE zones (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE racks (
    id SERIAL PRIMARY KEY,
    zone_id INTEGER NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL
);

CREATE TABLE shelves (
    id SERIAL PRIMARY KEY,
    rack_id INTEGER NOT NULL REFERENCES racks(id) ON DELETE CASCADE,
    level INTEGER NOT NULL
);

-- Products & stock
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    unit VARCHAR(20) NOT NULL
);

CREATE TABLE stock_items (
    id SERIAL PRIMARY KEY,
    shelf_id INTEGER NOT NULL REFERENCES shelves(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (shelf_id, product_id)
);

CREATE TABLE stock_movements (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    from_shelf_id INTEGER REFERENCES shelves(id) ON DELETE SET NULL,
    to_shelf_id INTEGER REFERENCES shelves(id) ON DELETE SET NULL,
    quantity INTEGER NOT NULL,
    moved_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);