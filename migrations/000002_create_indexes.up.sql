-- users
CREATE INDEX IF NOT EXISTS users_not_deleted_idx
    ON users (id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS users_role_idx
    ON users (role)
    WHERE deleted_at IS NULL;

-- locations
CREATE UNIQUE INDEX IF NOT EXISTS locations_address_room_uniq
    ON locations (address, room)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS locations_address_idx
    ON locations (address)
    WHERE deleted_at IS NULL;

-- equipment_categories
CREATE INDEX IF NOT EXISTS equipment_categories_not_deleted_idx
    ON equipment_categories (name)
    WHERE deleted_at IS NULL;

-- equipments
CREATE INDEX IF NOT EXISTS equipments_category_idx
    ON equipments (category_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS equipments_location_idx
    ON equipments (location_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS equipments_state_idx
    ON equipments (state)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS equipments_category_state_idx
    ON equipments (category_id, state)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS equipments_not_deleted_idx
    ON equipments (id)
    WHERE deleted_at IS NULL;