-- 009_governance.sql

CREATE TABLE permissions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level       INT NOT NULL CHECK (level BETWEEN 1 AND 4),
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    behavior    TEXT NOT NULL DEFAULT 'notify',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE principles (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name             TEXT NOT NULL UNIQUE,
    description      TEXT NOT NULL,
    evaluation_logic JSONB NOT NULL DEFAULT '{}',
    priority         INT NOT NULL DEFAULT 0,
    is_active        BOOLEAN NOT NULL DEFAULT true,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE control_rules (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    principle_id        UUID REFERENCES principles(id) ON DELETE CASCADE,
    target_entity_type  TEXT NOT NULL,
    target_entity_id    UUID,
    condition           JSONB NOT NULL DEFAULT '{}',
    action              TEXT NOT NULL,
    priority            INT NOT NULL DEFAULT 0,
    is_active           BOOLEAN NOT NULL DEFAULT true,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_principles_active ON principles(is_active);
CREATE INDEX idx_control_principle ON control_rules(principle_id);
CREATE INDEX idx_control_target ON control_rules(target_entity_type, target_entity_id);
