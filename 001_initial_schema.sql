-- 1. Users Table (role_id removed)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    department VARCHAR(100),
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Roles Table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);

-- 3. Permissions Table (Resource:Action)
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT
);

-- 4. User-Roles Junction Table (Many-to-Many)
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id INT REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- 5. Role-Permissions Junction Table (Many-to-Many)
CREATE TABLE role_permissions (
    role_id INT REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INT REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- --- SEEDING BASELINE DATA ---

-- Seed Roles
INSERT INTO roles (id, name, description) VALUES
(1, 'Administrator', 'Full system access'),
(2, 'Lab Manager', 'Manages hardware and checkouts'),
(3, 'Researcher', 'Standard hardware requester');

-- Seed Baseline Permissions
INSERT INTO permissions (id, name, description) VALUES
(1, 'inventory:read', 'Can view the hardware catalog'),
(2, 'inventory:create', 'Can add new hardware to the system'),
(3, 'checkout:create', 'Can request hardware'),
(4, 'checkout:approve', 'Can approve hardware requests');

-- Map Permissions to Roles (The power of RBAC)
-- Lab Managers get everything
INSERT INTO role_permissions (role_id, permission_id) VALUES 
(2, 1), (2, 2), (2, 3), (2, 4);

-- Researchers only get to read inventory and create checkouts
INSERT INTO role_permissions (role_id, permission_id) VALUES 
(3, 1), (3, 3);