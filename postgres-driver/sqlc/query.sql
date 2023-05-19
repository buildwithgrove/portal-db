-- name: CheckIDExists :one
SELECT EXISTS (
        SELECT 1
        FROM (
                SELECT id
                FROM portal_applications
                WHERE portal_applications.id = @id::VARCHAR
                UNION ALL
                SELECT id
                FROM aats
                WHERE aats.id = @id::VARCHAR
                UNION ALL
                SELECT id
                FROM accounts
                WHERE accounts.id = @id::VARCHAR
                UNION ALL
                SELECT id
                FROM users
                WHERE users.id = @id::VARCHAR
            ) AS id_table
    );
-- name: SelectPlans :many
SELECT plan_type,
    chain_ids,
    monthly_relay_limit,
    throughput_limit,
    application_limit,
    created_at,
    daily_limit
FROM pay_plans;
-- name: SelectGigastakeApplications :many
SELECT ga.id,
    ga.aat_id,
    ga.name,
    ga.chain_id,
    ga.chain_alias,
    ga.created_at,
    ga.updated_at,
    ga.deleted,
    a.gigastake,
    a.address,
    a.public_key,
    a.client_public_key,
    a.signature,
    a.private_key,
    a.version
FROM gigastake_applications AS ga
    JOIN aats AS a ON ga.aat_id = a.id;
-- name: SelectPortalApplications :many
WITH aats_agg AS (
    SELECT aats.portal_application_id,
        json_object_agg(
            aats.id,
            json_build_object(
                'address',
                aats.address,
                'public_key',
                aats.public_key,
                'private_key',
                aats.private_key,
                'client_public_key',
                aats.client_public_key,
                'signature',
                aats.signature,
                'version',
                aats.version
            )
        ) AS aats
    FROM aats aats
    GROUP BY aats.portal_application_id
),
notifications_agg AS (
    SELECT pn.application_id,
        json_object_agg(
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
                    WHERE event IS NOT NULL
                )
            )
        ) AS notifications
    FROM portal_application_notifications pn
    GROUP BY pn.application_id
),
whitelists_agg AS (
    SELECT paw.application_id,
        json_agg(
            json_build_object(
                'type',
                paw.type,
                'value',
                paw.value,
                'chain_id',
                paw.chain_id
            )
        ) AS whitelists
    FROM portal_application_whitelists paw
    GROUP BY paw.application_id
)
SELECT p.*,
    pas.secret_key,
    pas.secret_key_required,
    pas.monthly_relay_limit,
    pas.environment,
    COALESCE(aats_agg.aats, '[]'::json) AS aats,
    COALESCE(notifications_agg.notifications, '[]'::json) AS notifications,
    COALESCE(whitelists_agg.whitelists, '[]'::json) AS whitelists
FROM portal_applications p
    LEFT JOIN portal_application_settings pas ON p.id = pas.application_id
    LEFT JOIN aats_agg ON p.id = aats_agg.portal_application_id
    LEFT JOIN notifications_agg ON p.id = notifications_agg.application_id
    LEFT JOIN whitelists_agg ON p.id = whitelists_agg.application_id
WHERE (
        $1::BOOLEAN
        OR p.deleted = false
    );
-- name: InsertPortalApplication :one
INSERT INTO portal_applications (
        id,
        account_id,
        name,
        created_at,
        updated_at,
        request_timeout,
        first_date_surpassed,
        -- legacy field
        plan_type,
        -- legacy field
        daily_limit,
        -- legacy field
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
        $10
    )
RETURNING *;
-- name: InsertPortalApplicationAAT :one
INSERT INTO aats (
        id,
        portal_application_id,
        gigastake,
        address,
        public_key,
        private_key,
        client_public_key,
        signature,
        version
    )
VALUES ($1, $2, false, $3, $4, $5, $6, $7, $8)
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
-- name: GetUserEmail :one
SELECT email
FROM users
WHERE id = $1;
-- name: GetAccountOwnerEmail :one
SELECT users.email
FROM users
    JOIN account_user_access ON users.id = account_user_access.user_id
WHERE account_user_access.role_name = 'OWNER'
    AND account_user_access.account_id = $1;
-- name: UpdatePortalAppFields :exec
UPDATE portal_applications
SET name = COALESCE(NULLIF(@name::VARCHAR, ''), name),
    plan_type = COALESCE(NULLIF(@plan_type::VARCHAR, ''), plan_type),
    daily_limit = COALESCE($2, daily_limit),
    custom_limit = COALESCE($3, custom_limit),
    updated_at = $4
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
WITH user_auth_agg AS (
    SELECT user_id,
        json_object_agg(type, provider_user_id) AS provider_user_ids
    FROM user_auth_providers
    GROUP BY user_id
),
app_role_agg AS (
    SELECT user_id,
        account_id,
        json_object_agg(portal_application_id, role_name) AS portal_application_roles
    FROM account_user_access
    WHERE portal_application_id IS NOT NULL
    GROUP BY user_id,
        account_id
)
SELECT a.*,
    ai.covalent_api_key_free,
    ai.covalent_api_key_paid,
    json_agg(
        json_build_object(
            'user_id',
            u.id,
            'email',
            u.email,
            'owner',
            CASE
                WHEN au.role_name = 'OWNER' THEN true
                ELSE false
            END,
            'accepted',
            au.accepted,
            'provider_user_ids',
            uaa.provider_user_ids,
            'portal_application_roles',
            ara.portal_application_roles
        )
    ) AS users
FROM accounts AS a
    LEFT JOIN account_user_access AS au ON a.id = au.account_id
    LEFT JOIN account_integrations AS ai ON a.id = ai.account_id
    LEFT JOIN users AS u ON au.user_id = u.id
    LEFT JOIN user_roles AS ur ON au.role_name = ur.role_name
    LEFT JOIN user_auth_agg AS uaa ON u.id = uaa.user_id
    LEFT JOIN app_role_agg AS ara ON u.id = ara.user_id
    AND a.id = ara.account_id
WHERE (
        $1::BOOLEAN
        OR a.deleted = false
    )
GROUP BY a.id,
    ai.covalent_api_key_free,
    ai.covalent_api_key_paid;
-- name: SelectAccount :one
WITH user_auth_agg AS (
    SELECT user_id,
        json_object_agg(type, provider_user_id) AS provider_user_ids
    FROM user_auth_providers
    GROUP BY user_id
)
SELECT a.*,
    ai.covalent_api_key_free,
    ai.covalent_api_key_paid,
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
            uaa.provider_user_ids
        )
    ) AS users
FROM accounts AS a
    LEFT JOIN account_user_access AS au ON a.id = au.account_id
    LEFT JOIN account_integrations AS ai ON a.id = ai.account_id
    LEFT JOIN users AS u ON au.user_id = u.id
    LEFT JOIN user_roles AS ur ON au.role_name = ur.role_name
    LEFT JOIN user_auth_agg AS uaa ON u.id = uaa.user_id
WHERE a.id = $1
GROUP BY a.id,
    ai.covalent_api_key_free,
    ai.covalent_api_key_paid;
-- name: UpdateAccountFields :exec
UPDATE accounts
SET plan_type = COALESCE($1, plan_type),
    partner_chain_ids = COALESCE($2, partner_chain_ids),
    partner_throughput_limit = COALESCE($3, partner_throughput_limit),
    partner_application_limit = COALESCE($4, partner_application_limit),
    updated_at = $5
WHERE id = $6;
-- name: UpsertAccountIntegrations :one
INSERT INTO account_integrations (
        account_id,
        covalent_api_key_free,
        covalent_api_key_paid,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (account_id) DO
UPDATE
SET covalent_api_key_free = CASE
        WHEN EXCLUDED.covalent_api_key_free IS NOT NULL THEN EXCLUDED.covalent_api_key_free
        ELSE account_integrations.covalent_api_key_free
    END,
    covalent_api_key_paid = CASE
        WHEN EXCLUDED.covalent_api_key_paid IS NOT NULL THEN EXCLUDED.covalent_api_key_paid
        ELSE account_integrations.covalent_api_key_paid
    END,
    created_at = CASE
        WHEN EXCLUDED.created_at IS NOT NULL THEN EXCLUDED.created_at
        ELSE account_integrations.created_at
    END,
    updated_at = EXCLUDED.updated_at
RETURNING *;
-- name: SelectUserPermissions :many
SELECT aua.account_id,
    aua.user_id,
    aua.role_name,
    ur.permissions as permissions
FROM account_user_access as aua
    LEFT JOIN user_roles AS ur ON aua.role_name = ur.role_name
WHERE aua.accepted = true
    AND aua.user_id IS NOT NULL;
-- name: SelectUserIDs :many
SELECT user_id,
    provider_user_id
FROM user_auth_providers;
-- name: InsertAccount :one
INSERT INTO accounts (
        id,
        plan_type,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4)
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
-- name: CheckRedirectExists :one
SELECT EXISTS(
        SELECT 1
        FROM chain_gigastake_redirects
        WHERE chain_id = $1
            AND portal_application_id = $2
            AND domain = $3
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
-- name: CheckPortalAppExists :one
SELECT EXISTS(
        SELECT 1
        FROM portal_applications
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
WITH updated_user AS (
    UPDATE users
    SET email = $7
    WHERE id = $2
    RETURNING id,
        email
)
INSERT INTO account_user_access (
        account_id,
        user_id,
        portal_application_id,
        role_name,
        accepted,
        created_at,
        updated_at
    )
VALUES ($1, $2, '', $3, $4, $5, $6)
RETURNING account_user_access.user_id,
    account_user_access.role_name,
    account_user_access.accepted,
    COALESCE(
        (
            SELECT email
            FROM updated_user
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
            id,
            email,
            signed_up,
            created_at,
            updated_at
        )
    VALUES ($1, $2, false, $3, $4)
    RETURNING id,
        email
)
INSERT INTO account_user_access (
        account_id,
        user_id,
        portal_application_id,
        role_name,
        accepted,
        created_at,
        updated_at
    )
VALUES (
        $5,
        (
            SELECT id
            FROM inserted_user
        ),
        '',
        $6,
        false,
        $3,
        $4
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
    INSERT INTO users (id, email, signed_up, created_at, updated_at)
    VALUES ($1, $2, true, $3, $4) ON CONFLICT (email) DO NOTHING
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
                        WHERE users.email = $2
                    )
                )
        ),
        $5,
        $6,
        $7,
        $8
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
        request_timeout,
        log_limit_blocks,
        chain_aliases,
        allowed_methods,
        gigastake_redirect_domains,
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
    request_timeout = COALESCE(EXCLUDED.request_timeout, chains.request_timeout),
    log_limit_blocks = COALESCE(
        EXCLUDED.log_limit_blocks,
        chains.log_limit_blocks
    ),
    chain_aliases = COALESCE(EXCLUDED.chain_aliases, chains.chain_aliases),
    allowed_methods = COALESCE(EXCLUDED.allowed_methods, chains.allowed_methods),
    gigastake_redirect_domains = COALESCE(
        EXCLUDED.gigastake_redirect_domains,
        chains.gigastake_redirect_domains
    ),
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
        portal_application_id,
        alias,
        domain,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (chain_id, portal_application_id, domain) DO
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
    AND portal_application_id NOT IN (
        SELECT unnest(@portal_application_ids::VARCHAR [])
    );
-- name: DeleteGigastakeRedirect :exec
DELETE FROM chain_gigastake_redirects
WHERE chain_id = $1
    AND portal_application_id = $2
    and domain = $3;
-- name: UpsertChainCheck :exec
INSERT INTO chain_checks (
        chain_id,
        type,
        payload,
        result_key,
        allowance,
        evm_chain_id,
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
        $8
    ) ON CONFLICT (chain_id, type) DO
UPDATE
SET payload = COALESCE(EXCLUDED.payload, chain_checks.payload),
    result_key = COALESCE(EXCLUDED.result_key, chain_checks.result_key),
    allowance = COALESCE(EXCLUDED.allowance, chain_checks.allowance),
    evm_chain_id = COALESCE(EXCLUDED.evm_chain_id, chain_checks.evm_chain_id),
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
-- name: SelectGlobalBlockedContract :many
SELECT id,
    blocked_address
FROM global_blocked_contracts
WHERE active = true;
-- name: AddGlobalBlockedContract :exec
INSERT INTO global_blocked_contracts (blocked_address, created_at, updated_at)
VALUES ($1, $2, $3);
-- name: SetGlobalBlockedContractActive :one
UPDATE global_blocked_contracts
SET active = $2,
    updated_at = $3
WHERE blocked_address = $1
RETURNING id;
-- name: RemoveGlobalBlockedContract :one
DELETE FROM global_blocked_contracts
WHERE blocked_address = $1
RETURNING id;
