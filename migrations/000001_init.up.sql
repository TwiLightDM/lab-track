CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    surname    VARCHAR(100) NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(60)  NOT NULL,
    role       VARCHAR(20)  NOT NULL CHECK (role IN ('admin', 'lab_assistant', 'lab_manager', 'repairman')),
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE locations
(
    id         SERIAL PRIMARY KEY,
    address    VARCHAR(255) NOT NULL,
    room       INT          NOT NULL,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE equipment_categories
(
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(100) NOT NULL UNIQUE,
    descriptions TEXT,
    created_at   TIMESTAMP,
    updated_at   TIMESTAMP,
    deleted_at   TIMESTAMP
);

CREATE TABLE equipments
(
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    inventory_number VARCHAR(100) NOT NULL UNIQUE,
    serial_number    VARCHAR(100) UNIQUE,
    state            VARCHAR(30)  NOT NULL
        CHECK (state IN (
                         'idle',
                         'in_use',
                         'maintenance_required',
                         'in_repair',
                         'written_off'
            )),
    description      TEXT,
    category_id      INTEGER      NOT NULL REFERENCES equipment_categories (id) ON DELETE RESTRICT,
    location_id      INTEGER      NOT NULL REFERENCES locations (id) ON DELETE RESTRICT,
    image_url        TEXT,
    created_at       TIMESTAMP,
    updated_at       TIMESTAMP,
    deleted_at       TIMESTAMP
);