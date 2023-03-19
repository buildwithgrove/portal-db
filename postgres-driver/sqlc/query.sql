-- name: SelectPortalApplications :many
SELECT p.*,
    paa.address,
    paa.public_key,
    paa.private_key,
    paa.client_public_key,
    paa.signature,
    paa.version,
    pas.secret_key,
    pas.secret_key_required,
    pas.monthly_relay_limit,
    pas.environment,
    -- legacy field
    pso.duration,
    -- legacy field
    pso.sticky_max,
    -- legacy field
    pso.stickiness,
    -- legacy field
    pso.origins,
    COALESCE(
        (
            SELECT json_object_agg(
                    pn.type,
                    json_build_object(
                        'active',
                        pn.active,
                        'destination',
                        pn.destination,
                        'trigger',
                        pn.trigger,
                        'events',
                        (
                            SELECT json_object_agg(
                                    event,
                                    true
                                )
                            FROM (
                                    SELECT unnest(pn.events) AS event
                                ) subquery
                        )
                    )
                )
            FROM portal_application_notifications pn
            WHERE pn.application_id = p.id
        ),
        '[]'::json
    )::json AS notifications,
    COALESCE(
        (
            SELECT json_agg(
                    json_build_object(
                        'type',
                        paw.type,
                        'value',
                        paw.value,
                        'chain_id',
                        paw.chain_id
                    )
                )
            FROM portal_application_whitelists paw
            WHERE paw.application_id = p.id
        ),
        '[]'::json
    )::json AS whitelists
FROM portal_applications p
    LEFT JOIN portal_application_aats paa ON p.id = paa.application_id
    LEFT JOIN portal_application_settings pas ON p.id = pas.application_id -- legacy table
    LEFT JOIN stickiness_options pso ON p.id = pso.lb_id
WHERE (
        @include_deleted::BOOLEAN
        OR p.deleted = false
    )
GROUP BY p.id,
    paa.address,
    paa.public_key,
    paa.private_key,
    paa.client_public_key,
    paa.signature,
    paa.version,
    pas.secret_key,
    pas.secret_key_required,
    pas.monthly_relay_limit,
    pas.environment,
    pso.duration,
    pso.sticky_max,
    pso.stickiness,
    pso.origins;
-- name: InsertPortalApplication :one
INSERT INTO portal_applications (
        id,
        account_id,
        name,
        gigastake,
        staked,
        created_at,
        updated_at,
        application_ids,
        request_timeout,
        gigastake_redirect,
        first_date_surpassed,
        custom_limit
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12
    )
RETURNING *;
-- name: InsertStickinessOption :one
INSERT INTO stickiness_options (
        lb_id,
        duration,
        sticky_max,
        stickiness,
        origins
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: InsertPortalApplicationAAT :one
INSERT INTO portal_application_aats (
        application_id,
        address,
        public_key,
        private_key,
        client_public_key,
        signature,
        version
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
-- name: InsertPortalApplicationSetting :one
INSERT INTO portal_application_settings (
        application_id,
        secret_key,
        secret_key_required,
        monthly_relay_limit,
        environment
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: UpdatePortalAppName :exec
UPDATE portal_applications
SET name = $2,
    updated_at = $3
WHERE id = $1;
-- name: UpdatePortalAppSettings :exec
UPDATE portal_application_settings
SET secret_key = COALESCE($2, secret_key),
    secret_key_required = COALESCE($3, secret_key_required),
    monthly_relay_limit = COALESCE($4, monthly_relay_limit),
    environment = COALESCE($5, environment),
    favorited_chain_ids = COALESCE($6, favorited_chain_ids),
    updated_at = COALESCE($7, updated_at)
WHERE application_id = $1;
-- name: UpdateUpsertPortalAppNotification :exec
INSERT INTO portal_application_notifications (
        application_id,
        type,
        active,
        destination,
        trigger,
        events,
        updated_at
    )
SELECT $1,
    @type::notification_type,
    @active::BOOLEAN,
    @destination::VARCHAR(255),
    @trigger::VARCHAR(255),
    @events::notification_event [],
    $2
WHERE @active IS true ON CONFLICT (application_id, type) DO
UPDATE
SET active = EXCLUDED.active,
    destination = EXCLUDED.destination,
    trigger = EXCLUDED.trigger,
    events = EXCLUDED.events,
    updated_at = EXCLUDED.updated_at
WHERE EXCLUDED.active IS true;
-- name: UpdateDeletePortalAppNotification :exec
DELETE FROM portal_application_notifications
WHERE application_id = $1
    and type = $2;
-- name: UpdateInsertWhitelists :exec
INSERT INTO portal_application_whitelists (
        application_id,
        type,
        chain_id,
        value,
        created_at
    )
VALUES(
        $1,
        unnest(@types::whitelist_type []),
        NULLIF(unnest(@chain_ids::VARCHAR []), ''),
        unnest(@values::VARCHAR []),
        @created_at::TIMESTAMPTZ
    ) ON CONFLICT (application_id, chain_id, type, value) DO NOTHING;
-- name: UpdateDeleteWhitelists :exec
DELETE FROM portal_application_whitelists
WHERE (
        type IN ('methods', 'contracts')
        AND application_id = $1
        AND NOT EXISTS (
            SELECT 1
            FROM unnest(@types::whitelist_type []) t1(type),
                unnest(@values::VARCHAR []) t2(value),
                unnest(@chain_ids::VARCHAR []) t3(chain_id)
            WHERE t1.type = portal_application_whitelists.type
                AND t2.value = portal_application_whitelists.value
                AND t3.chain_id = portal_application_whitelists.chain_id
        )
    )
    OR (
        type IN ('blockchains', 'origins', 'userAgents')
        AND application_id = $1
        AND NOT EXISTS (
            SELECT 1
            FROM unnest(@types::whitelist_type []) t1(type),
                unnest(@values::VARCHAR []) t2(value)
            WHERE t1.type = portal_application_whitelists.type
                AND t2.value = portal_application_whitelists.value
        )
    );
-- name: UpdateFirstDatesSurpassed :exec
UPDATE portal_applications
SET first_date_surpassed = @first_date_surpassed
WHERE id = ANY (@application_ids::VARCHAR []);
-- name: DeletePortalApp :exec
UPDATE portal_applications
SET deleted = true,
    deleted_at = $2
WHERE id = $1;
-- name: SelectAccounts :many
SELECT a.*,
    p.chain_ids,
    p.monthly_relay_limit,
    p.throughput_limit,
    p.application_limit,
    json_agg(
        json_build_object(
            'user_id',
            u.id,
            'email',
            u.email,
            'role_name',
            ur.role_name,
            'accepted',
            au.accepted,
            'provider_user_ids',
            (
                SELECT json_object_agg(type, provider_user_id)
                FROM user_auth_providers
                WHERE user_id = u.id
            )
        )
    ) AS users,
    -- legacy field
    p.daily_limit
FROM accounts AS a
    LEFT JOIN account_user_access AS au ON a.id = au.account_id
    LEFT JOIN pay_plans AS p ON a.plan_type = p.plan_type
    LEFT JOIN users AS u ON au.user_id = u.id
    LEFT JOIN user_roles AS ur ON au.role_name = ur.role_name
WHERE (
        @include_deleted::BOOLEAN
        OR a.deleted = false
    )
GROUP BY a.id,
    p.plan_type,
    p.chain_ids,
    p.monthly_relay_limit,
    p.throughput_limit,
    p.application_limit,
    p.daily_limit;
-- name: InsertAccount :one
INSERT INTO accounts (
        plan_type,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3)
RETURNING *;
-- name: DeleteAccount :exec
UPDATE accounts
SET deleted = true,
    deleted_at = $2
WHERE id = $1;
-- name: CheckChainExists :one
SELECT EXISTS(
        SELECT 1
        FROM chains
        WHERE id = $1
            AND deleted = false
    );
-- name: CheckPlanTypeExists :one
SELECT EXISTS(
        SELECT 1
        FROM pay_plans
        WHERE plan_type = $1
    );
-- name: CheckUserEmail :one
SELECT email
FROM users
WHERE id = $1;
-- name: CheckUserIDFromEmail :one
SELECT id
FROM users
WHERE email = $1;
-- name: CheckAccountExists :one
SELECT EXISTS(
        SELECT 1
        FROM accounts
        WHERE id = $1
            AND deleted = false
    );
-- name: CheckUserExists :one
SELECT EXISTS(
        SELECT 1
        FROM users
        WHERE id = $1
    );
-- name: CheckAccountUserExists :one
SELECT EXISTS (
        SELECT 1
        FROM account_user_access
        WHERE user_id = $1
            AND account_id = $2
    );
-- name: CheckAccountUserRole :one
SELECT role_name
FROM account_user_access
WHERE user_id = $1
    AND account_id = $2;
-- name: CheckAccountUserAccepted :one
SELECT accepted
FROM account_user_access
WHERE user_id = $1
    AND account_id = $2;
-- name: InsertAccountUserAccess :one
INSERT INTO account_user_access (
        account_id,
        user_id,
        role_name,
        accepted,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING account_user_access.user_id,
    account_user_access.role_name,
    account_user_access.accepted,
    COALESCE(
        (
            SELECT email
            FROM users
            WHERE id = $2
        ),
        ''
    )::VARCHAR(320) AS user_email,
    (
        SELECT json_object_agg(type, provider_user_id)
        FROM user_auth_providers
        WHERE user_id = account_user_access.user_id
    ) as provider_user_ids;
-- name: InsertAccountUserAccessNoUser :one
WITH inserted_user AS (
    INSERT INTO users (
            email,
            signed_up,
            created_at,
            updated_at
        )
    VALUES ($1, false, $4, $5)
    RETURNING id,
        email
)
INSERT INTO account_user_access (
        account_id,
        user_id,
        role_name,
        accepted,
        created_at,
        updated_at
    )
VALUES (
        $2,
        (
            SELECT id
            FROM inserted_user
        ),
        $3,
        false,
        $4,
        $5
    )
RETURNING account_user_access.user_id,
    account_user_access.role_name,
    account_user_access.accepted,
    (
        SELECT email
        FROM inserted_user
    ) AS user_email;
-- name: UpdateAccountUserRole :exec
UPDATE account_user_access
SET role_name = $3,
    updated_at = $4
WHERE account_id = $1
    AND user_id = $2;
-- name: UpdateAccountOwnerToAdmin :exec
UPDATE account_user_access
SET role_name = 'ADMIN'
WHERE account_user_access.account_id = $1
    AND role_name = 'OWNER'
    AND user_id = (
        SELECT user_id
        FROM account_user_access
        WHERE account_id = $1
            AND role_name = 'OWNER'
    );
-- name: CreateUserNewSignUp :one
WITH inserted_user AS (
    INSERT INTO users (email, signed_up, created_at, updated_at)
    VALUES ($1, true, $2, $3) ON CONFLICT (email) DO NOTHING
    RETURNING id
)
INSERT INTO user_auth_providers (
        user_id,
        type,
        provider,
        provider_user_id,
        federated
    )
VALUES (
        (
            SELECT COALESCE(
                    (
                        SELECT id
                        FROM inserted_user
                    ),
                    (
                        SELECT id
                        FROM users
                        WHERE users.email = $1
                    )
                )
        ),
        $4,
        $5,
        $6,
        $7
    )
RETURNING user_id;
-- name: UpdateUserAcceptedInvite :exec
WITH inserted_provider AS (
    INSERT INTO user_auth_providers (
            user_id,
            type,
            provider,
            provider_user_id,
            federated
        )
    VALUES ($1, $2, $3, $4, $5)
    RETURNING user_id
),
updated_access AS (
    UPDATE account_user_access
    SET accepted = true
    WHERE user_id = (
            SELECT user_id
            FROM inserted_provider
        )
        AND account_id = $6
)
UPDATE users
SET signed_up = true
WHERE id = (
        SELECT user_id
        FROM inserted_provider
    );
-- name: DeleteAccountUser :exec
DELETE FROM account_user_access
WHERE account_id = $1
    AND user_id = $2;
-- name: GetPortalUserID :one
SELECT user_id
FROM user_auth_providers
WHERE provider_user_id = $1;
-- name: GetUserDataFromPortalUserID :one
SELECT users.*,
    json_agg(user_auth_providers.*) AS auth_providers
FROM users
    LEFT JOIN user_auth_providers ON users.id = user_auth_providers.user_id
WHERE users.id = $1
GROUP BY users.id;
-- name: GetUserDataFromAuthProviderUserID :one
SELECT users.*,
    json_agg(user_auth_providers.*) AS auth_providers
FROM users
    LEFT JOIN user_auth_providers ON users.id = user_auth_providers.user_id
WHERE user_auth_providers.provider_user_id = $1
GROUP BY users.id;
-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id;
-- name: SelectChains :many
SELECT c.*,
    COALESCE(
        json_agg(DISTINCT ca) FILTER (
            WHERE ca.id IS NOT NULL
        ),
        '[]'
    )::json AS chain_altruists,
    COALESCE(
        json_agg(DISTINCT cgr) FILTER (
            WHERE cgr.id IS NOT NULL
        ),
        '[]'
    )::json AS chain_gigastake_redirects,
    COALESCE(
        json_agg(DISTINCT cc) FILTER (
            WHERE cc.id IS NOT NULL
        ),
        '[]'
    )::json AS chain_checks
FROM chains c
    LEFT JOIN chain_altruists ca ON c.id = ca.chain_id
    LEFT JOIN chain_gigastake_redirects cgr ON c.id = cgr.chain_id
    LEFT JOIN chain_checks cc ON c.id = cc.chain_id
WHERE (
        @include_deleted::BOOLEAN
        OR c.deleted = false
    )
GROUP BY c.id;
-- name: UpsertChain :one
INSERT INTO chains (
        id,
        blockchain,
        description,
        enforce_result,
        path,
        ticker,
        blockchain_id,
        request_timeout,
        log_limit_blocks,
        chain_aliases,
        allowed_methods,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13
    ) ON CONFLICT (id) DO
UPDATE
SET blockchain = COALESCE(EXCLUDED.blockchain, chains.blockchain),
    description = COALESCE(EXCLUDED.description, chains.description),
    enforce_result = COALESCE(EXCLUDED.enforce_result, chains.enforce_result),
    path = COALESCE(EXCLUDED.path, chains.path),
    ticker = COALESCE(EXCLUDED.ticker, chains.ticker),
    blockchain_id = COALESCE(EXCLUDED.blockchain_id, chains.blockchain_id),
    request_timeout = COALESCE(EXCLUDED.request_timeout, chains.request_timeout),
    log_limit_blocks = COALESCE(
        EXCLUDED.log_limit_blocks,
        chains.log_limit_blocks
    ),
    chain_aliases = COALESCE(EXCLUDED.chain_aliases, chains.chain_aliases),
    allowed_methods = COALESCE(EXCLUDED.allowed_methods, chains.allowed_methods),
    updated_at = EXCLUDED.updated_at
RETURNING id;
-- name: UpsertChainAltruist :exec
INSERT INTO chain_altruists (
        chain_id,
        url,
        auth,
        auth_type,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (chain_id, url) DO
UPDATE
SET auth = COALESCE(EXCLUDED.auth, chain_altruists.auth),
    auth_type = COALESCE(EXCLUDED.auth_type, chain_altruists.auth_type),
    updated_at = EXCLUDED.updated_at;
-- name: DeleteUnusedChainAltruists :exec
DELETE FROM chain_altruists
WHERE chain_id = $1
    AND url NOT IN (
        SELECT unnest(@urls::VARCHAR [])
    );
-- name: UpsertChainGigastakeRedirect :exec
INSERT INTO chain_gigastake_redirects (
        chain_id,
        account_id,
        alias,
        domain,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (account_id) DO
UPDATE
SET chain_id = COALESCE(
        EXCLUDED.chain_id,
        chain_gigastake_redirects.chain_id
    ),
    alias = COALESCE(EXCLUDED.alias, chain_gigastake_redirects.alias),
    domain = COALESCE(
        EXCLUDED.domain,
        chain_gigastake_redirects.domain
    ),
    updated_at = EXCLUDED.updated_at;
-- name: DeleteUnusedChainGigastakeRedirects :exec
DELETE FROM chain_gigastake_redirects
WHERE chain_id = $1
    AND account_id NOT IN (
        SELECT unnest(@account_ids::INTEGER [])
    );
-- name: UpsertChainCheck :exec
INSERT INTO chain_checks (
        chain_id,
        type,
        payload,
        result_key,
        allowance,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7
    ) ON CONFLICT (chain_id, type) DO
UPDATE
SET payload = COALESCE(EXCLUDED.payload, chain_checks.payload),
    result_key = COALESCE(EXCLUDED.result_key, chain_checks.result_key),
    allowance = COALESCE(EXCLUDED.allowance, chain_checks.allowance),
    updated_at = EXCLUDED.updated_at;
-- name: DeleteUnusedChainChecks :exec
DELETE FROM chain_checks
WHERE chain_id = $1
    AND type NOT IN (
        SELECT unnest(@types::chain_check_type [])
    );
-- name: UpdateChainActive :one
UPDATE chains
SET active = $2,
    updated_at = $3
WHERE id = $1
RETURNING active;
