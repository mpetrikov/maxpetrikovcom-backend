CREATE TABLE lab_sessions (
    id UUID PRIMARY KEY,
    lab_id UUID NOT NULL,
    user_id UUID NOT NULL,

    status TEXT NOT NULL DEFAULT 'pending',

    namespace TEXT NULL,
    pod_name TEXT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NULL,

    failure_reason TEXT NULL,

    CONSTRAINT fk_lab_sessions_lab
        FOREIGN KEY (lab_id)
        REFERENCES labs(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_lab_sessions_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_lab_sessions_status
        CHECK (
            status IN (
                'pending',
                'provisioning',
                'running',
                'stopping',
                'stopped',
                'expired',
                'failed'
            )
        )
);

CREATE INDEX idx_lab_sessions_user_id
    ON lab_sessions(user_id);

CREATE INDEX idx_lab_sessions_lab_id
    ON lab_sessions(lab_id);

CREATE INDEX idx_lab_sessions_status
    ON lab_sessions(status);

CREATE INDEX idx_lab_sessions_expires_at
    ON lab_sessions(expires_at);