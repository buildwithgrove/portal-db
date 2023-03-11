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
    LEFT JOIN portal_application_settings pas ON p.id = pas.application_id
    LEFT JOIN stickiness_options pso ON p.id = pso.application_id
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
        application_id,
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
-- name: UpdateInsertChainWhitelists :exec
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
        unnest(@chain_ids::VARCHAR []),
        unnest(@values::VARCHAR []),
        @created_at::TIMESTAMPTZ
    ) ON CONFLICT (application_id, chain_id, type, value) DO NOTHING;
-- name: UpdateDeletePortalAppWhitelists :exec
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
