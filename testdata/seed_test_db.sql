INSERT INTO pay_plans (
        plan_type,
        monthly_relay_limit,
        throughput_limit,
        application_limit
    )
VALUES (
        'test_plan',
        4200000,
        2000,
        3
    );
INSERT INTO accounts (plan_type)
VALUES ('test_plan');
INSERT INTO portal_applications (
        id,
        account_id,
        name,
        gigastake,
        staked,
        application_id,
        request_timeout,
        gigastake_redirect,
        first_date_surpassed,
        created_at,
        updated_at
    )
VALUES (
        'test_app_3487u329rfn23f9',
        1,
        'pokt_app_123',
        true,
        false,
        'test_app_47hfnths73j2se7',
        5000,
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2308rj09r23r9r',
        1,
        'pokt_app_456',
        false,
        true,
        'test_app_43jr9304urj30fj',
        10000,
        false,
        '2022-11-11 11:11:11.000000',
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
        'test_app_3487u329rfn23f9',
        60,
        300,
        true,
        '{ "chrome-extension://", "moz-extension://" }'
    ),
    (
        'test_app_2308rj09r23r9r',
        30,
        600,
        true,
        '{ "https://example.com", "https://test.com" }'
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
        'test_app_3487u329rfn23f9',
        'test_34715cae753e67c75fbb340442e7de8e',
        'test_34715cae753e67c75fbb340442e7de8e',
        'test_11b8d394ca331d7c7a71ca1896d630f6',
        'test_89a3af6a587aec02cfade6f5000424c2',
        'test_1dc39a2e5a84a35bf030969a0b3231f7',
        '0.0.1',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2308rj09r23r9r',
        'test_8237c72345f12d1b1a8b64a1a7f66fa4',
        'test_8237c72345f12d1b1a8b64a1a7f66fa4',
        'test_2e83c836a29b423a47d8e18c779fd422',
        'test_04c71d90a92f40416b6f1d7d8af17e02',
        'test_f48d33b30ddaf60a1e5bb50d2ba8da5a',
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
        'test_app_3487u329rfn23f9',
        'test_40f482d91a5ef2300ebb4e2308c',
        true,
        2500000,
        'production',
        ARRAY ['0001', '0053'],
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2308rj09r23r9r',
        'test_9c9e3b193cfba5348f93bb2f3e3fb794',
        false,
        1500000,
        'production',
        ARRAY ['0021', '0064'],
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
        'test_app_3487u329rfn23f9',
        'blockchains',
        '0053',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_3487u329rfn23f9',
        'contracts',
        '0x1234567890abcdef',
        '0001',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_3487u329rfn23f9',
        'methods',
        'GET',
        '0001',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_3487u329rfn23f9',
        'origins',
        'https://example.com',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_3487u329rfn23f9',
        'userAgents',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2308rj09r23r9r',
        'blockchains',
        '0021',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2308rj09r23r9r',
        'contracts',
        '0x0987654321abcdef',
        '0064',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2308rj09r23r9r',
        'methods',
        'POST',
        '0064',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2308rj09r23r9r',
        'origins',
        'https://test.com',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2308rj09r23r9r',
        'userAgents',
        'Mozilla/5.0 (Linux; Android 10; SM-A205U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36',
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
        'test_app_3487u329rfn23f9',
        true,
        'email',
        'test@test.com',
        'trigger123',
        ARRAY ['quarter'::notification_event,'threeQuarters'::notification_event,'full'::notification_event],
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2308rj09r23r9r',
        true,
        'email',
        'email@pokt.network',
        'trigger456',
        ARRAY ['full'::notification_event,'half'::notification_event],
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
