-- 008_verification.sql

CREATE TABLE verification_reports (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id         UUID REFERENCES workflow_instances(id) ON DELETE SET NULL,
    task_id             UUID REFERENCES tasks(id) ON DELETE SET NULL,
    result_score        DOUBLE PRECISION,
    path_score          DOUBLE PRECISION,
    environment_score   DOUBLE PRECISION,
    overall_score       DOUBLE PRECISION,
    conclusion          TEXT NOT NULL DEFAULT '',
    suggestions         JSONB NOT NULL DEFAULT '[]',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE review_assignments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id       UUID NOT NULL REFERENCES verification_reports(id) ON DELETE CASCADE,
    level           TEXT NOT NULL,
    reviewer_id     UUID,
    reviewer_type   TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending',
    result          JSONB,
    completed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_verif_workflow ON verification_reports(workflow_id);
CREATE INDEX idx_verif_task ON verification_reports(task_id);
CREATE INDEX idx_review_report ON review_assignments(report_id);
CREATE INDEX idx_review_reviewer ON review_assignments(reviewer_id);
