CREATE TABLE session.email_verification_code (
    id            UUID          NOT NULL DEFAULT gen_random_uuid(),
    user_id       UUID          NOT NULL,
    code_hash     TEXT          NOT NULL,
    expires_at    TIMESTAMPTZ   NOT NULL DEFAULT (now() + INTERVAL '15 minutes'),
    is_used       BOOLEAN       NOT NULL DEFAULT FALSE,
    requested_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    ip_address    VARCHAR(45),

    CONSTRAINT pk_email_verification_code
        PRIMARY KEY (id)
);

CREATE INDEX ix_email_verification_user_id
    ON session.email_verification_code (user_id);