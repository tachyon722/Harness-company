-- 007_observability.sql

CREATE TABLE traces (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id     UUID REFERENCES workflow_instances(id) ON DELETE SET NULL,
    status          TEXT NOT NULL DEFAULT 'active',
    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ,
    metadata        JSONB NOT NULL DEFAULT '{}'
);

CREATE TABLE spans (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trace_id        UUID NOT NULL REFERENCES traces(id) ON DELETE CASCADE,
    parent_span_id  UUID REFERENCES spans(id) ON DELETE SET NULL,
    span_type       TEXT NOT NULL,
    entity_id       UUID,
    entity_type     TEXT,
    actor_id        UUID,
    actor_type      TEXT,
    input           JSONB,
    output          JSONB,
    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ,
    duration_ms     INT,
    metadata        JSONB NOT NULL DEFAULT '{}'
);

CREATE TABLE metrics (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    metric_type     TEXT NOT NULL,
    metric_name     TEXT NOT NULL,
    entity_id       UUID,
    entity_type     TEXT,
    value           DOUBLE PRECISION NOT NULL,
    recorded_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata        JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_spans_trace_id ON spans(trace_id);
CREATE INDEX idx_spans_actor ON spans(actor_id);
CREATE INDEX idx_spans_type ON spans(span_type);
CREATE INDEX idx_metrics_type ON metrics(metric_type, metric_name);
CREATE INDEX idx_metrics_recorded_at ON metrics(recorded_at);
