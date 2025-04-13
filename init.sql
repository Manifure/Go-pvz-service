CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     email TEXT UNIQUE NOT NULL,
                                     password TEXT NOT NULL,
                                     role TEXT CHECK (role IN ('client', 'moderator', 'employee')) NOT NULL
);

CREATE TABLE pvz (
                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                     city TEXT NOT NULL,
                     registered_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS acceptances (
                                           id UUID PRIMARY KEY,
                                           created_at TIMESTAMP NOT NULL DEFAULT now(),
                                           pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
                                           status TEXT NOT NULL CHECK (status IN ('in_progress', 'closed'))
);

CREATE TABLE IF NOT EXISTS items (
                                     id UUID PRIMARY KEY,
                                     received_at TIMESTAMP NOT NULL DEFAULT now(),
                                     type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь'))
);

CREATE TABLE acceptance_items (
                                  acceptance_id UUID REFERENCES acceptances(id),
                                  item_id UUID REFERENCES items(id),
                                  PRIMARY KEY (acceptance_id, item_id)
);
