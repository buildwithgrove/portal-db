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
-- name: CheckUserEmail :one
SELECT email
FROM users
WHERE id = $1;
-- name: CheckUserIDFromEmail :one
SELECT id
FROM users
WHERE email = $1;
-- name: CheckUserExists :one
SELECT EXISTS(
        SELECT 1
        FROM users
        WHERE id = $1
    );
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
-- name: CreateUserNewSignUp :one
WITH inserted_user AS (
    INSERT INTO users (email, signed_up, created_at, updated_at)
    VALUES ($1, true, $2, $3)
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
            SELECT id
            FROM inserted_user
        ),
        $4,
        $5,
        $6,
        $7
    )
RETURNING (
        SELECT id
        FROM inserted_user
    ) as user_id;
-- name: CreateUserProviderSignedUp :one
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
)
UPDATE users
SET signed_up = true
WHERE id = (
        SELECT user_id
        FROM inserted_provider
    )
RETURNING (
        SELECT user_id
        FROM inserted_provider
    ) as user_id;
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
