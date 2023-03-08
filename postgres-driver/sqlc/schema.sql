-- Enums
CREATE TYPE auth_providers AS ENUM ('auth0');
CREATE TYPE auth_sign_in AS ENUM ('github', 'username');
CREATE TYPE auth_type AS ENUM ('basic_auth', 'bearer_token', 'none');
CREATE TYPE chain_check_type AS ENUM ('archival', 'chain', 'merge', 'sync');
CREATE TYPE environment AS ENUM ('production', 'test');
CREATE TYPE notification_event AS ENUM (
    'full',
    'half',
    'quarter',
    'signedUp',
    'threeQuarters'
);
CREATE TYPE notification_type AS ENUM ('email', 'portal', 'webhook');
CREATE TYPE permissions_enum AS ENUM (
    'read:endpoint',
    'write:endpoint',
    'delete:endpoint',
    'transfer:endpoint'
);
CREATE TYPE whitelist_type AS ENUM (
    'blockchains',
    'contracts',
    'methods',
    'origins',
    'userAgents'
);
-- Blockchains Tables
CREATE TABLE blockchains (
    id VARCHAR(4) NOT NULL,
    blockchain VARCHAR(100) NOT NULL,
    description VARCHAR(100) NOT NULL,
    enforce_result VARCHAR(4) NOT NULL,
    path VARCHAR(100) NOT NULL,
    ticker VARCHAR(20) NOT NULL,
    chain_id INT,
    request_timeout INT,
    log_limit_blocks INT,
    blockchain_aliases VARCHAR(100) ARRAY,
    allowed_methods VARCHAR(10) ARRAY,
    active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted BOOLEAN DEFAULT false,
    PRIMARY KEY (id)
);
CREATE TABLE blockchain_altruists (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    blockchain_id VARCHAR(4) NOT NULL,
    url VARCHAR(255) NOT NULL,
    auth VARCHAR(100),
    auth_type auth_type,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (blockchain_id, url),
    CONSTRAINT blockchain_altruists_blockchain_id_fk FOREIGN KEY (blockchain_id) REFERENCES blockchains(id) ON DELETE CASCADE
);
CREATE TABLE blockchain_gigastake_redirects (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    account_id BIGINT NOT NULL UNIQUE,
    blockchain_id VARCHAR(4) NOT NULL,
    alias VARCHAR(100) NOT NULL,
    domain VARCHAR(100) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT blockchain_gigastake_redirects_blockchain_id_fk FOREIGN KEY (blockchain_id) REFERENCES blockchains(id) ON DELETE CASCADE,
    CONSTRAINT blockchain_gigastake_redirects_account_id_fk FOREIGN KEY (account_id) REFERENCES accounts(id)
);
CREATE TABLE blockchain_checks (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    blockchain_id VARCHAR(4) NOT NULL,
    type chain_check_type,
    payload VARCHAR(255) NOT NULL,
    result_key VARCHAR(100),
    allowance INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (blockchain_id, type),
    CONSTRAINT blockchain_checks_blockchain_id_fk FOREIGN KEY (blockchain_id) REFERENCES blockchains(id) ON DELETE CASCADE,
    CONSTRAINT sync_allowance_check CHECK (
        (
            type = 'sync'
            AND allowance IS NOT NULL
        )
        OR (
            type <> 'sync'
            AND allowance IS NULL
        )
    )
);
-- Portal Application Tables
CREATE TABLE portal_applications (
    id VARCHAR(24) NOT NULL,
    account_id BIGINT,
    name VARCHAR(255) NOT NULL,
    gigastake BOOLEAN NOT NULL,
    -- legacy field
    application_id VARCHAR(24),
    -- legacy field
    request_timeout INT,
    -- legacy field
    gigastake_redirect BOOLEAN,
    -- legacy field
    first_date_surpassed TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted BOOLEAN DEFAULT false,
    PRIMARY KEY (id),
    CONSTRAINT portal_application__account_id_fk FOREIGN KEY(account_id) REFERENCES accounts(id)
);
-- legacy table
CREATE TABLE IF NOT EXISTS stickiness_options (
    id INT GENERATED ALWAYS AS IDENTITY,
    application_id VARCHAR(24) NOT NULL UNIQUE,
    duration TEXT,
    sticky_max INT,
    stickiness BOOLEAN,
    origins VARCHAR ARRAY,
    PRIMARY KEY (id),
    CONSTRAINT stickiness_options_app_id_fk FOREIGN KEY(application_id) REFERENCES portal_applications(id) ON DELETE CASCADE
);
CREATE TABLE portal_application_aats (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    application_id VARCHAR(24) NOT NULL UNIQUE,
    address VARCHAR(40) NOT NULL,
    public_key VARCHAR(64) NOT NULL,
    private_key VARCHAR(400) NOT NULL,
    client_public_key VARCHAR(64) NOT NULL,
    signature VARCHAR(128) NOT NULL,
    version VARCHAR(10) NOT NULL,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT portal_aats_app_id_fk FOREIGN KEY (application_id) REFERENCES portal_applications(id) ON DELETE CASCADE
);
CREATE TABLE portal_application_settings (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    application_id VARCHAR(24) NOT NULL UNIQUE,
    secret_key VARCHAR(64) NOT NULL,
    secret_key_required BOOLEAN NOT NULL,
    monthly_relay_limit INT NOT NULL,
    environment environment NOT NULL,
    favorited_blockchain_ids VARCHAR(4) ARRAY REFERENCES blockchains (id),
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT portal_settings_app_id_fk FOREIGN KEY (application_id) REFERENCES portal_applications(id) ON DELETE CASCADE
);
CREATE TABLE portal_application_whitelists (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    application_id VARCHAR(24) NOT NULL,
    type whitelist_type NOT NULL,
    value VARCHAR(255) NOT NULL,
    blockchain_id VARCHAR(4),
    created_at TIMESTAMP NULL,
    UNIQUE (application_id, value, type),
    UNIQUE (application_id, value, type, blockchain_id),
    PRIMARY KEY (id),
    CONSTRAINT portal_whitelists_app_id_fk FOREIGN KEY (application_id) REFERENCES portal_applications(id) ON DELETE CASCADE,
    CONSTRAINT check_blockchain_id_for_methods_contracts CHECK (
        (
            type NOT IN ('methods', 'contracts')
            AND blockchain_id IS NULL
        )
        OR (
            type IN ('methods', 'contracts')
            AND blockchain_id IS NOT NULL
        )
    )
);
CREATE TABLE portal_application_notifications (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    application_id VARCHAR(24) NOT NULL UNIQUE,
    active BOOLEAN NOT NULL,
    type notification_type NOT NULL,
    destination VARCHAR(255),
    trigger VARCHAR(255),
    events notification_event ARRAY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (application_id, type),
    CONSTRAINT portal_notifications_app_id_fk FOREIGN KEY (application_id) REFERENCES portal_applications(id) ON DELETE CASCADE
);
-- Users Tables
CREATE TABLE users (
    id VARCHAR(320),
    email VARCHAR(320) NOT NULL UNIQUE,
    auth_provider auth_providers,
    sign_in_type auth_sign_in,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);
CREATE TABLE user_roles (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    role_name VARCHAR(25) UNIQUE,
    permissions permissions_enum ARRAY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);
-- Plans Tables
CREATE TABLE pay_plans (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    plan_type VARCHAR(25) UNIQUE,
    blockchain_ids VARCHAR(4) ARRAY REFERENCES blockchains (id),
    monthly_relay_limit INT NOT NULL,
    throughput_limit INT NOT NULL,
    application_limit INT NOT NULL,
    -- legacy field
    daily_limit INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);
-- Accounts Tables
CREATE TABLE accounts (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    plan_type VARCHAR(25) NOT NULL,
    partner_blockchain_ids VARCHAR(4) ARRAY REFERENCES blockchains (id),
    partner_throughput_limit INT,
    partner_application_limit INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted BOOLEAN DEFAULT false,
    PRIMARY KEY (id),
    CONSTRAINT accounts_pay_plans_fk FOREIGN KEY (plan_type) REFERENCES pay_plans(plan_type)
);
CREATE TABLE account_user_access (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    account_id BIGINT NOT NULL,
    user_id VARCHAR(320) NOT NULL,
    role_name VARCHAR(25) NOT NULL,
    accepted BOOLEAN,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT account_user_access_account_id_fk FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    CONSTRAINT account_user_access_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT account_user_access_user_roles_fk FOREIGN KEY (role_name) REFERENCES user_roles(role_name),
    CONSTRAINT account_user_access_unique_account_user UNIQUE (account_id, user_id)
);
-- Blocked Contracts Tables
CREATE TABLE global_blocked_contracts (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    blocked_address VARCHAR(255) UNIQUE,
    active BOOLEAN DEFAULT true,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);
-- Listener Notification Function
CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$
DECLARE data json;
notification json;
BEGIN -- Convert the old or new row to JSON, based on the kind of action.
-- Action = DELETE?             -> OLD row
-- Action = INSERT or UPDATE?   -> NEW row
IF (TG_OP = 'DELETE') THEN data = row_to_json(OLD);
ELSE data = row_to_json(NEW);
END IF;
-- Contruct the notification as a JSON string.
notification = json_build_object(
    'table',
    TG_TABLE_NAME,
    'action',
    TG_OP,
    'data',
    data
);
-- Execute pg_notify(channel, notification)
PERFORM pg_notify('events', notification::text);
-- Result is ignored since this is an AFTER trigger
RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- Listener Notification Triggers
CREATE TRIGGER portal_applications_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON portal_applications FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER portal_application_aats_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON portal_application_aats FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER portal_application_settings_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON portal_application_settings FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER portal_application_whitelists_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON portal_application_whitelists FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER portal_application_notifications_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON portal_application_notifications FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER accounts_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON accounts FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER account_user_access_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON account_user_access FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER users_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON users FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER user_roles_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON user_roles FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER pay_plans_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON pay_plans FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER blockchains_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON blockchains FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER blockchain_altruists_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON blockchain_altruists FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER blockchain_gigastake_redirects_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON blockchain_gigastake_redirects FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER blockchain_checks_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON blockchain_checks FOR EACH ROW EXECUTE PROCEDURE notify_event();
CREATE TRIGGER global_blocked_contracts_notify_event
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON global_blocked_contracts FOR EACH ROW EXECUTE PROCEDURE notify_event();
