-- name: SelectBlockchains :many
SELECT b.blockchain_id,
    b.altruist,
    b.blockchain,
    b.blockchain_aliases,
    b.chain_id,
    b.chain_id_check,
    b.description,
    b.enforce_result,
    b.log_limit_blocks,
    b.network,
    b.path,
    b.request_timeout,
    b.ticker,
    b.active,
    s.synccheck AS s_sync_check,
    s.allowance AS s_allowance,
    s.body AS s_body,
    s.path AS s_path,
    s.result_key AS s_result_key,
    COALESCE(redirects.r, '[]') AS redirects,
    b.created_at,
    b.updated_at
FROM blockchains AS b
    LEFT JOIN sync_check_options AS s ON b.blockchain_id = s.blockchain_id
    LEFT JOIN LATERAL (
        SELECT json_agg(
                json_build_object(
                    'alias',
                    r.alias,
                    'loadBalancerID',
                    r.loadbalancer,
                    'domain',
                    r.domain
                )
            ) AS r
        FROM redirects AS r
        WHERE b.blockchain_id = r.blockchain_id
    ) redirects ON true
ORDER BY b.blockchain_id ASC;
-- name: SelectPayPlans :many
SELECT plan_type,
    daily_limit
FROM pay_plans
ORDER BY plan_type ASC;
-- name: InsertBlockchain :exec
INSERT into blockchains (
        blockchain_id,
        active,
        altruist,
        blockchain,
        blockchain_aliases,
        chain_id,
        chain_id_check,
        description,
        enforce_result,
        log_limit_blocks,
        network,
        path,
        request_timeout,
        ticker,
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
        $13,
        $14,
        $15,
        $16
    );
-- name: InsertRedirect :exec
INSERT into redirects (
        blockchain_id,
        alias,
        loadbalancer,
        domain,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    );
-- name: InsertSyncCheckOptions :exec
INSERT into sync_check_options (
        blockchain_id,
        synccheck,
        allowance,
        body,
        path,
        result_key
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    );
-- name: ActivateBlockchain :exec
UPDATE blockchains
SET active = $2,
    updated_at = $3
WHERE blockchain_id = $1;
-- name: SelectApplications :many
WITH app_whitelists AS (
    SELECT application_id
    FROM whitelist_contracts
    UNION
    SELECT application_id
    FROM whitelist_methods
)
SELECT a.application_id,
    a.contact_email,
    a.created_at,
    a.description,
    a.dummy,
    a.name,
    a.owner,
    a.status,
    a.updated_at,
    a.url,
    a.user_id,
    a.first_date_surpassed,
    ga.address AS ga_address,
    ga.client_public_key AS ga_client_public_key,
    ga.private_key AS ga_private_key,
    ga.public_key AS ga_public_key,
    ga.signature AS ga_signature,
    ga.version AS ga_version,
    gs.secret_key,
    gs.secret_key_required,
    gs.whitelist_blockchains,
    gs.whitelist_origins,
    gs.whitelist_user_agents,
    ns.signed_up,
    ns.on_quarter,
    ns.on_half,
    ns.on_three_quarters,
    ns.on_full,
    al.custom_limit,
    al.pay_plan,
    pp.daily_limit as plan_limit,
    CASE
        WHEN wc.application_id IS NOT NULL THEN json_agg(
            json_build_object(
                'blockchain_id',
                wc.blockchain_id,
                'contracts',
                wc.contracts
            )
        )::VARCHAR
        ELSE null
    END as whitelist_contracts,
    CASE
        WHEN wm.application_id IS NOT NULL THEN json_agg(
            json_build_object(
                'blockchain_id',
                wm.blockchain_id,
                'methods',
                wm.methods
            )
        )::VARCHAR
        ELSE null
    END as whitelist_methods
FROM applications AS a
    LEFT JOIN gateway_aat AS ga ON a.application_id = ga.application_id
    LEFT JOIN gateway_settings AS gs ON a.application_id = gs.application_id
    LEFT JOIN notification_settings AS ns ON a.application_id = ns.application_id
    LEFT JOIN app_limits AS al ON a.application_id = al.application_id
    LEFT JOIN pay_plans AS pp ON al.pay_plan = pp.plan_type
    LEFT JOIN whitelist_contracts wc ON a.application_id = wc.application_id
    LEFT JOIN whitelist_methods wm ON a.application_id = wm.application_id
GROUP BY a.application_id,
    a.contact_email,
    a.created_at,
    a.description,
    a.dummy,
    a.name,
    a.owner,
    a.status,
    a.updated_at,
    a.url,
    a.user_id,
    a.first_date_surpassed,
    ga.address,
    ga.client_public_key,
    ga.private_key,
    ga.public_key,
    ga.signature,
    ga.version,
    gs.secret_key,
    gs.secret_key_required,
    gs.whitelist_blockchains,
    gs.whitelist_origins,
    gs.whitelist_user_agents,
    ns.signed_up,
    ns.on_quarter,
    ns.on_half,
    ns.on_three_quarters,
    ns.on_full,
    al.custom_limit,
    al.pay_plan,
    pp.daily_limit,
    wc.application_id,
    wm.application_id;
-- name: SelectOneApplication :one
WITH app_whitelists AS (
    SELECT application_id
    FROM whitelist_contracts
    UNION
    SELECT application_id
    FROM whitelist_methods
)
SELECT a.application_id,
    a.contact_email,
    a.created_at,
    a.description,
    a.dummy,
    a.name,
    a.owner,
    a.status,
    a.updated_at,
    a.url,
    a.user_id,
    a.first_date_surpassed,
    ga.address AS ga_address,
    ga.client_public_key AS ga_client_public_key,
    ga.private_key AS ga_private_key,
    ga.public_key AS ga_public_key,
    ga.signature AS ga_signature,
    ga.version AS ga_version,
    gs.secret_key,
    gs.secret_key_required,
    gs.whitelist_blockchains,
    gs.whitelist_origins,
    gs.whitelist_user_agents,
    ns.signed_up,
    ns.on_quarter,
    ns.on_half,
    ns.on_three_quarters,
    ns.on_full,
    al.custom_limit,
    al.pay_plan,
    pp.daily_limit as plan_limit,
    CASE
        WHEN wc.application_id IS NOT NULL THEN json_agg(
            json_build_object(
                'blockchain_id',
                wc.blockchain_id,
                'contracts',
                wc.contracts
            )
        )::VARCHAR
        ELSE null
    END as whitelist_contracts,
    CASE
        WHEN wm.application_id IS NOT NULL THEN json_agg(
            json_build_object(
                'blockchain_id',
                wm.blockchain_id,
                'methods',
                wm.methods
            )
        )::VARCHAR
        ELSE null
    END as whitelist_methods
FROM applications AS a
    LEFT JOIN gateway_aat AS ga ON a.application_id = ga.application_id
    LEFT JOIN gateway_settings AS gs ON a.application_id = gs.application_id
    LEFT JOIN notification_settings AS ns ON a.application_id = ns.application_id
    LEFT JOIN app_limits AS al ON a.application_id = al.application_id
    LEFT JOIN pay_plans AS pp ON al.pay_plan = pp.plan_type
    LEFT JOIN whitelist_contracts wc ON a.application_id = wc.application_id
    LEFT JOIN whitelist_methods wm ON a.application_id = wm.application_id
WHERE a.application_id = $1
GROUP BY a.application_id,
    a.contact_email,
    a.created_at,
    a.description,
    a.dummy,
    a.name,
    a.owner,
    a.status,
    a.updated_at,
    a.url,
    a.user_id,
    a.first_date_surpassed,
    ga.address,
    ga.client_public_key,
    ga.private_key,
    ga.public_key,
    ga.signature,
    ga.version,
    gs.secret_key,
    gs.secret_key_required,
    gs.whitelist_blockchains,
    gs.whitelist_origins,
    gs.whitelist_user_agents,
    ns.signed_up,
    ns.on_quarter,
    ns.on_half,
    ns.on_three_quarters,
    ns.on_full,
    al.custom_limit,
    al.pay_plan,
    pp.daily_limit,
    wc.application_id,
    wm.application_id;
-- name: SelectAppLimit :one
SELECT application_id,
    pay_plan,
    custom_limit
FROM app_limits
WHERE application_id = $1;
-- name: SelectGatewaySettings :one
SELECT gs.application_id AS application_id,
    gs.secret_key AS secret_key,
    gs.secret_key_required AS secret_key_required,
    gs.whitelist_blockchains AS whitelist_blockchains,
    json_agg(
        json_build_object(
            'blockchain_id',
            wc.blockchain_id,
            'contracts',
            wc.contracts
        )
    )::VARCHAR as whitelist_contracts,
    json_agg(
        json_build_object(
            'blockchain_id',
            wm.blockchain_id,
            'methods',
            wm.methods
        )
    )::VARCHAR as whitelist_methods,
    gs.whitelist_origins AS whitelist_origins,
    gs.whitelist_user_agents AS whitelist_user_agents
FROM gateway_settings AS gs
    LEFT JOIN whitelist_contracts AS wc ON gs.application_id = wc.application_id
    LEFT JOIN whitelist_methods AS wm ON gs.application_id = wm.application_id
WHERE gs.application_id = $1;
-- name: SelectNotificationSettings :one
SELECT application_id,
    signed_up,
    on_quarter,
    on_half,
    on_three_quarters,
    on_full
FROM notification_settings
WHERE application_id = $1;
-- name: InsertApplication :exec
INSERT into applications (
        application_id,
        user_id,
        name,
        contact_email,
        description,
        owner,
        url,
        status,
        dummy,
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
        $11
    );
-- name: InsertAppLimit :exec
INSERT into app_limits (application_id, pay_plan, custom_limit)
VALUES ($1, $2, $3);
-- name: InsertGatewayAAT :exec
INSERT into gateway_aat (
        application_id,
        address,
        client_public_key,
        private_key,
        public_key,
        signature,
        version
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7
    );
-- name: InsertGatewaySettings :exec
INSERT into gateway_settings (
        application_id,
        secret_key,
        secret_key_required
    )
VALUES (
        $1,
        $2,
        $3
    );
-- name: InsertNotificationSettings :exec
INSERT into notification_settings (
        application_id,
        signed_up,
        on_quarter,
        on_half,
        on_three_quarters,
        on_full
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    );
-- name: UpsertApplication :exec
INSERT INTO applications AS a (
        application_id,
        name,
        status,
        first_date_surpassed,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (application_id) DO
UPDATE
SET name = COALESCE(EXCLUDED.name, a.name),
    status = COALESCE(EXCLUDED.status, a.status),
    first_date_surpassed = COALESCE(
        EXCLUDED.first_date_surpassed,
        a.first_date_surpassed
    );
-- name: UpsertAppLimit :exec
INSERT INTO app_limits AS al (
        application_id,
        pay_plan,
        custom_limit
    )
VALUES ($1, $2, $3) ON CONFLICT (application_id) DO
UPDATE
SET pay_plan = COALESCE(EXCLUDED.pay_plan, al.pay_plan),
    custom_limit = COALESCE(EXCLUDED.custom_limit, al.custom_limit);
-- name: UpsertGatewaySettings :exec
INSERT INTO gateway_settings AS gs (
        application_id,
        secret_key,
        secret_key_required,
        whitelist_origins,
        whitelist_user_agents,
        whitelist_blockchains
    )
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (application_id) DO
UPDATE
SET secret_key = COALESCE(EXCLUDED.secret_key, gs.secret_key),
    secret_key_required = COALESCE(
        EXCLUDED.secret_key_required,
        gs.secret_key_required
    ),
    whitelist_origins = COALESCE(EXCLUDED.whitelist_origins, gs.whitelist_origins),
    whitelist_user_agents = COALESCE(
        EXCLUDED.whitelist_user_agents,
        gs.whitelist_user_agents
    ),
    whitelist_blockchains = COALESCE(
        EXCLUDED.whitelist_blockchains,
        gs.whitelist_blockchains
    );
-- name: UpsertWhitelistContracts :exec
WITH data (application_id, blockchain_id, contracts) AS (
    VALUES (
            @application_id::VARCHAR,
            @blockchain_id::VARCHAR,
            @contracts::VARCHAR []
        )
)
INSERT INTO whitelist_contracts (application_id, blockchain_id, contracts)
SELECT application_id,
    blockchain_id,
    contracts
FROM data ON CONFLICT (application_id, blockchain_id) DO
UPDATE
SET contracts = excluded.contracts;
-- name: UpsertWhitelistMethods :exec
WITH data (application_id, blockchain_id, methods) AS (
    VALUES (
            @application_id::VARCHAR,
            @blockchain_id::VARCHAR,
            @methods::VARCHAR []
        )
)
INSERT INTO whitelist_methods (application_id, blockchain_id, methods)
SELECT application_id,
    blockchain_id,
    methods::VARCHAR []
FROM data ON CONFLICT (application_id, blockchain_id) DO
UPDATE
SET methods = EXCLUDED.methods;
-- name: UpsertNotificationSettings :exec
INSERT INTO notification_settings AS ns (
        application_id,
        signed_up,
        on_quarter,
        on_half,
        on_three_quarters,
        on_full
    )
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (application_id) DO
UPDATE
SET signed_up = COALESCE(EXCLUDED.signed_up, ns.signed_up),
    on_quarter = COALESCE(EXCLUDED.on_quarter, ns.on_quarter),
    on_half = COALESCE(EXCLUDED.on_half, ns.on_half),
    on_three_quarters = COALESCE(EXCLUDED.on_three_quarters, ns.on_three_quarters),
    on_full = COALESCE(EXCLUDED.on_full, ns.on_full);
-- name: UpdateFirstDateSurpassed :exec
UPDATE applications
SET first_date_surpassed = @first_date_surpassed
WHERE application_id = ANY (@application_ids::VARCHAR []);
-- name: RemoveApp :exec
UPDATE applications
SET status = COALESCE($2, status)
WHERE application_id = $1;
-- name: SelectLoadBalancers :many
SELECT lb.lb_id,
    lb.name,
    lb.request_timeout,
    lb.gigastake,
    lb.gigastake_redirect,
    lb.user_id,
    so.duration AS s_duration,
    so.sticky_max AS s_sticky_max,
    so.stickiness AS s_stickiness,
    so.origins AS s_origins,
    STRING_AGG(la.app_id, ',') AS app_ids,
    COALESCE(user_access.ua, '[]') AS users,
    lb.created_at,
    lb.updated_at
FROM loadbalancers AS lb
    LEFT JOIN stickiness_options AS so ON lb.lb_id = so.lb_id
    LEFT JOIN lb_apps AS la ON lb.lb_id = la.lb_id
    LEFT JOIN LATERAL (
        SELECT jsonb_agg(
                json_build_object(
                    'userID',
                    ua.user_id,
                    'roleName',
                    ua.role_name,
                    'email',
                    ua.email,
                    'accepted',
                    ua.accepted
                )
            ) AS ua
        FROM user_access AS ua
        WHERE lb.lb_id = ua.lb_id
    ) user_access ON true
GROUP BY lb.lb_id,
    lb.lb_id,
    lb.name,
    lb.created_at,
    lb.updated_at,
    lb.request_timeout,
    lb.gigastake,
    lb.gigastake_redirect,
    lb.user_id,
    so.duration,
    so.sticky_max,
    so.stickiness,
    so.origins,
    user_access.ua
ORDER BY lb.lb_id ASC;
-- name: SelectOneLoadBalancer :one
SELECT lb.lb_id,
    lb.name,
    lb.request_timeout,
    lb.gigastake,
    lb.gigastake_redirect,
    lb.user_id,
    so.duration,
    so.sticky_max,
    so.stickiness,
    so.origins,
    STRING_AGG(la.app_id, ',') AS app_ids,
    COALESCE(user_access.ua, '[]') AS users,
    lb.created_at,
    lb.updated_at
FROM loadbalancers AS lb
    LEFT JOIN stickiness_options AS so ON lb.lb_id = so.lb_id
    LEFT JOIN lb_apps AS la ON lb.lb_id = la.lb_id
    LEFT JOIN LATERAL (
        SELECT jsonb_agg(
                json_build_object(
                    'userID',
                    ua.user_id,
                    'roleName',
                    ua.role_name,
                    'email',
                    ua.email,
                    'accepted',
                    ua.accepted
                )
            ) AS ua
        FROM user_access AS ua
        WHERE lb.lb_id = ua.lb_id
    ) user_access ON true
WHERE lb.lb_id = $1
GROUP BY lb.lb_id,
    lb.lb_id,
    lb.name,
    lb.created_at,
    lb.updated_at,
    lb.request_timeout,
    lb.gigastake,
    lb.gigastake_redirect,
    lb.user_id,
    so.duration,
    so.sticky_max,
    so.stickiness,
    so.origins,
    user_access.ua;
-- name: SelectUserRoles :many
SELECT ua.lb_id,
    ua.user_id,
    ur.permissions as permissions
FROM user_access as ua
    LEFT JOIN user_roles AS ur ON ua.role_name = ur.name;
-- name: InsertLoadBalancer :exec
INSERT into loadbalancers (
        lb_id,
        name,
        user_id,
        request_timeout,
        gigastake,
        gigastake_redirect,
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
    );
-- name: InsertStickinessOptions :exec
INSERT INTO stickiness_options (
        lb_id,
        duration,
        sticky_max,
        stickiness,
        origins
    )
VALUES ($1, $2, $3, $4, $5);
-- name: InsertUserAccess :exec
INSERT INTO user_access (
        lb_id,
        role_name,
        user_id,
        email,
        accepted,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7);
-- name: UpdateUserAccess :exec
UPDATE user_access as ua
SET role_name = COALESCE($3, ua.role_name),
    updated_at = $4
WHERE ua.user_id = $1
    AND ua.lb_id = $2;
-- name: DeleteUserAccess :exec
DELETE FROM user_access
WHERE user_id = $1
    AND lb_id = $2;
-- name: UpsertStickinessOptions :exec
INSERT INTO stickiness_options AS so (
        lb_id,
        duration,
        sticky_max,
        stickiness,
        origins
    )
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (lb_id) DO
UPDATE
SET duration = COALESCE(EXCLUDED.duration, so.duration),
    sticky_max = COALESCE(EXCLUDED.sticky_max, so.sticky_max),
    stickiness = COALESCE(EXCLUDED.stickiness, so.stickiness),
    origins = COALESCE(EXCLUDED.origins, so.origins);
-- name: InsertLbApps :exec
INSERT into lb_apps (lb_id, app_id)
SELECT @lb_id,
    unnest(@app_ids::VARCHAR []);
-- name: UpdateLB :exec
UPDATE loadbalancers AS l
SET name = COALESCE($2, l.name),
    updated_at = $3
WHERE l.lb_id = $1;
-- name: RemoveLB :exec
UPDATE loadbalancers
SET user_id = '',
    updated_at = $2
WHERE lb_id = $1;
