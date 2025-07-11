CREATE TABLE IF NOT EXISTS weather (
    city TEXT PRIMARY KEY,
    temperature INT
);

INSERT INTO weather (city, temperature) VALUES
('Saint-Petersburg', 20)
ON CONFLICT (city) DO NOTHING;