INSERT INTO portal_applications (
        id,
        name,
        gigastake,
        application_id,
        request_timeout,
        gigastake_redirect,
        first_date_surpassed,
        created_at,
        updated_at
    )
VALUES (
        'test_app_3487u329rfn23f',
        'pokt_app_123',
        true,
        'test_app_47hfnths73j2se',
        5000,
        true,
        '2022-01-01 00:00:00',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO stickiness_options (
        application_id,
        duration,
        sticky_max,
        stickiness,
        origins
    )
VALUES (
        'test_app_3487u329rfn23f',
        60,
        300,
        true,
        '{ "chrome-extension://", "moz-extension://" }'
    );
INSERT INTO portal_application_aats (
        application_id,
        address,
        public_key,
        private_key,
        client_public_key,
        signature,
        version,
        updated_at
    )
VALUES (
        'test_app_3487u329rfn23f',
        'test_34715cae753e67c75fbb340442e7de8e',
        'test_34715cae753e67c75fbb340442e7de8e',
        'test_11b8d394ca331d7c7a71ca1896d630f6',
        'test_89a3af6a587aec02cfade6f5000424c2',
        'test_1dc39a2e5a84a35bf030969a0b3231f7',
        '0.0.1',
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO portal_application_settings (
        application_id,
        secret_key,
        secret_key_required,
        monthly_relay_limit,
        environment,
        favorited_blockchain_ids,
        updated_at
    )
VALUES (
        'test_app_3487u329rfn23f',
        'test_40f482d91a5ef2300ebb4e2308c',
        true,
        2500000,
        'production',
        ARRAY ['0001', '0053'],
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO portal_application_whitelists (
        application_id,
        type,
        value,
        blockchain_id,
        created_at
    )
VALUES (
        'test_app_3487u329rfn23f',
        'blockchains',
        '0053',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_3487u329rfn23f',
        'contracts',
        '0x1234567890abcdef',
        '0001',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_3487u329rfn23f',
        'methods',
        'GET',
        '0001',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_3487u329rfn23f',
        'origins',
        'https://example.com',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_3487u329rfn23f',
        'userAgents',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',
        NULL,
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO portal_application_notifications (
        application_id,
        active,
        type,
        destination,
        trigger,
        events,
        created_at,
        updated_at
    )
VALUES (
        'test_app_3487u329rfn23f',
        true,
        'email',
        'test@test.com',
        'trigger123',
        ARRAY [
        'quarter'::notification_event,
        'threeQuarters'::notification_event,
        'full'::notification_event
    ],
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
