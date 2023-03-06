-- name: SelectPortalApplications :many
SELECT p.*,
    pso.duration,
    pso.sticky_max,
    pso.stickiness,
    pso.origins,
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
            pn.events
        )
    ) AS notifications,
    json_agg(
        json_build_object(
            'type',
            paw.type,
            'value',
            paw.value,
            'blockchain_id',
            paw.blockchain_id
        )
    ) AS whitelists
FROM portal_applications AS p
    LEFT JOIN stickiness_options AS pso ON p.id = pso.application_id
    LEFT JOIN portal_application_aats AS paa ON p.id = paa.application_id
    LEFT JOIN portal_application_settings AS pas ON p.id = pas.application_id
    LEFT JOIN portal_application_notifications AS pn ON p.id = pn.application_id
    LEFT JOIN portal_application_whitelists AS paw ON p.id = paw.application_id
GROUP BY p.id,
    pso.duration,
    pso.sticky_max,
    pso.stickiness,
    pso.origins,
    paa.address,
    paa.public_key,
    paa.private_key,
    paa.client_public_key,
    paa.signature,
    paa.version,
    pas.secret_key,
    pas.secret_key_required,
    pas.monthly_relay_limit,
    pas.environment;
