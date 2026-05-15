CREATE TABLE IF NOT EXISTS logs (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('processing', 'completed', 'error')) DEFAULT 'processing',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS nodes (
    id SERIAL PRIMARY KEY,
    log_id INT NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    guid VARCHAR(64) NOT NULL,
    description VARCHAR(255),
    node_type INT NOT NULL,
    system_image_guid VARCHAR(64),
    base_version INT,
    class_version INT
);

CREATE INDEX IF NOT EXISTS idx_nodes_guid ON nodes(guid);

CREATE TABLE IF NOT EXISTS ports (
    id SERIAL PRIMARY KEY,
    log_id INT NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    node_id INT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    guid VARCHAR(64),
    number INT NOT NULL,
    lid INT,
    state INT,
    physical_state INT,
    link_speed INT,
    link_width INT
);

CREATE TABLE IF NOT EXISTS nodes_info (
    id SERIAL PRIMARY KEY,
    node_id INT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    serial_number VARCHAR(255),
    part_number VARCHAR(255),
    revision VARCHAR(255),
    product_name VARCHAR(255)
);
