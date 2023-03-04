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
    'transfer:endpoint',
);
CREATE TYPE whitelist_type AS ENUM (
    'blockchains',
    'contracts',
    'methods',
    'origins',
    'userAgents'
);
-- Portal Application Tables
CREATE TABLE portal_applications (
    id VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    gigastake BOOLEAN NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);
CREATE TABLE portal_application_aats (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    portal_application_id VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    public_key VARCHAR NOT NULL,
    private_key VARCHAR NOT NULL,
    client_public_key VARCHAR NOT NULL,
    signature VARCHAR NOT NULL,
    version VARCHAR NOT NULL,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT portal_aats_portal_application_id_fk FOREIGN KEY (portal_application_id) REFERENCES portal_applications(id) ON DELETE CASCADE,
    CONSTRAINT portal_aats_unique_portal_application_id_fk UNIQUE (portal_application_id)
);
CREATE TABLE portal_application_settings (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    portal_application_id VARCHAR NOT NULL,
    secret_key VARCHAR NOT NULL,
    secret_key_required BOOLEAN NOT NULL,
    app_monthly_relay_limit INT NOT NULL,
    environment environment NOT NULL,
    favorited_blockchain_ids VARCHAR [],
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT portal_settings_portal_application_id_fk FOREIGN KEY (portal_application_id) REFERENCES portal_applications(id) ON DELETE CASCADE,
    CONSTRAINT portal_settings_unique_portal_application_id_fk UNIQUE (portal_application_id)
);
CREATE TABLE portal_application_whitelists (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    portal_application_id VARCHAR NOT NULL,
    type whitelist_type NOT NULL,
    value VARCHAR NOT NULL,
    blockchain_id VARCHAR,
    created_at TIMESTAMP NULL,
    UNIQUE (application_id, value, type),
    UNIQUE (application_id, value, type, blockchain_id),
    PRIMARY KEY (id),
    CONSTRAINT portal_whitelists_portal_application_id_fk FOREIGN KEY (portal_application_id) REFERENCES portal_applications(id) ON DELETE CASCADE,
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
    portal_application_id VARCHAR NOT NULL,
    active BOOLEAN NOT NULL,
    type notification_type NOT NULL,
    destination VARCHAR,
    trigger VARCHAR,
    events notification_event [],
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT portal_notifications_portal_application_id_fk FOREIGN KEY (portal_application_id) REFERENCES portal_applications(id) ON DELETE CASCADE
);
-- Accounts Tables
CREATE TABLE accounts (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    portal_application_id VARCHAR NOT NULL,
    plan_type VARCHAR UNIQUE NOT NULL,
    partner_blockchain_ids VARCHAR [],
    partner_throughput_limit INT,
    partner_application_limit INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT accounts_portal_application_id_fk FOREIGN KEY (portal_application_id) REFERENCES portal_applications(id),
    CONSTRAINT accounts_unique_portal_application_id_fk UNIQUE (portal_application_id),
    CONSTRAINT account_partner_blockchains_fk FOREIGN KEY (partner_blockchain_ids) REFERENCES blockchains(id),
    CONSTRAINT accounts_pay_plans_fk FOREIGN KEY (plan_type) REFERENCES pay_plans(plan_type)
);
CREATE TABLE account_user_access (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    account_id BIGINT NOT NULL,
    user_id VARCHAR NOT NULL,
    role_name VARCHAR NOT NULL,
    accepted BOOLEAN,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT account_user_access_account_id_fk FOREIGN KEY (account_id) REFERENCES accounts(id),
    CONSTRAINT account_user_access_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT account_user_access_unique_account_user FOREIGN KEY (account_id, user_id) UNIQUE,
    CONSTRAINT account_user_access_user_roles_fk FOREIGN KEY (role_name) REFERENCES user_roles(role_name)
);
-- Users Tables
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR NOT NULL UNIQUE,
    auth_provider auth_providers,
    sign_in_type auth_sign_in,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE TABLE user_roles (
    id VARCHAR NOT NULL,
    role_name VARCHAR,
    permissions permissions_enum [],
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);
-- Plans Tables
CREATE TABLE pay_plans (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    plan_type VARCHAR UNIQUE,
    blockchain_ids VARCHAR [],
    monthly_relay_limit INT NOT NULL,
    throughput_limit INT NOT NULL,
    application_limit INT NOT NULL,
    legacy_daily_limit INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT pay_plans_blockchains_fk FOREIGN KEY (blockchain_ids) REFERENCES blockchains(id)
);
-- Blockchains Tables
CREATE TABLE blockchains (
    id VARCHAR NOT NULL,
    blockchain VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    enforce_result VARCHAR NOT NULL,
    path VARCHAR NOT NULL,
    ticker VARCHAR NOT NULL,
    chain_id VARCHAR,
    request_timeout INT,
    log_limit_blocks INT,
    blockchain_aliases VARCHAR [],
    allowed_methods VARCHAR [],
    active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);
CREATE TABLE blockchain_altruists (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    blockchain_id VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
    auth VARCHAR,
    auth_type auth_type,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT blockchain_altruists_blockchain_id_fk FOREIGN KEY (blockchain_id) REFERENCES blockchains(id) ON DELETE CASCADE
);
CREATE TABLE blockchain_gigastake_redirects (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    blockchain_id VARCHAR NOT NULL,
    alias VARCHAR NOT NULL,
    -- TODO change this name?,
    protocol_app_id  VARCHAR NOT NULL,
    domain VARCHAR NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT blockchain_gigastake_redirects_blockchain_id_fk FOREIGN KEY (blockchain_id) REFERENCES blockchains(id) ON DELETE CASCADE
);
CREATE TABLE blockchain_checks (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    blockchain_id VARCHAR NOT NULL,
    type chain_check_type,
    payload VARCHAR NOT NULL,
    result_key VARCHAR,
    -- TODO enforce allowance if type = sync
    allowance INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT blockchain_checks_blockchain_id_fk FOREIGN KEY (blockchain_id) REFERENCES blockchains(id) ON DELETE CASCADE
);
-- Blocked Contracts Tables
CREATE TABLE global_blocked_contracts (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    blocked_address VARCHAR UNIQUE,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);
