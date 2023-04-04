-- Insert pay_plans
INSERT INTO pay_plans (
        plan_type,
        chain_ids,
        monthly_relay_limit,
        throughput_limit,
        application_limit,
        created_at,
        updated_at,
        daily_limit
    )
VALUES (
        'basic_plan',
        ARRAY ['0001', '0053'],
        5000000,
        5000,
        2,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000',
        1000
    ),
    (
        'pro_plan',
        ARRAY ['0001', '0053','0021', '0064'],
        10000000,
        10000,
        5,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000',
        5000
    ),
    (
        'enterprise_plan',
        ARRAY ['0001', '0053','0021', '0064', '0034'],
        20000000,
        20000,
        10,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000',
        10000
    ),
    (
        'developer_plan',
        ARRAY ['0001', '0053','0021', '0034'],
        500000,
        500,
        1,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000',
        100
    ),
    (
        'startup_plan',
        ARRAY ['0001', '0053', '0064', '0034'],
        1000000,
        1000,
        5,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000',
        500
    );
-- Insert accounts
INSERT INTO accounts (
        id,
        plan_type,
        partner_chain_ids,
        partner_throughput_limit,
        partner_application_limit,
        created_at,
        updated_at
    )
VALUES (
        'account_1',
        'basic_plan',
        ARRAY ['0001', '0053'],
        2000,
        1,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_2',
        'pro_plan',
        ARRAY ['0001', '0053','0021', '0064'],
        5000,
        3,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_3',
        'startup_plan',
        ARRAY ['0001', '0053', '0064', '0034'],
        1000,
        2,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_4',
        'enterprise_plan',
        ARRAY ['0001'],
        1000,
        2,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_5',
        'basic_plan',
        ARRAY ['0006', '0040'],
        6000,
        1,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
-- Insert users
INSERT INTO users (id, email, signed_up, created_at, updated_at)
VALUES (
        'user_1',
        'james.holden123@test.com',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_2',
        'paul.atreides456@test.com',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_3',
        'ellen.ripley789@test.com',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_4',
        'ulfric.stormcloak123@test.com',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_5',
        'chrisjen.avasarala1@test.com',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_6',
        'amos.burton789@test.com',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_7',
        'frodo.baggins123@test.com',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_8',
        'rick.deckard456@test.com',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_9',
        'tyrion.lannister789@test.com',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_10',
        'daenerys.targaryen123@test.com',
        false,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_11',
        'bernard.marx@test.com',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO user_auth_providers (
        user_id,
        provider,
        type,
        provider_user_id,
        federated,
        created_at
    )
VALUES (
        'user_1',
        'auth0',
        'auth0_username',
        'auth0|james_holden',
        false,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_1',
        'auth0',
        'auth0_github',
        'github|james_holden',
        true,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_2',
        'auth0',
        'auth0_username',
        'auth0|paul_atreides',
        false,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_2',
        'auth0',
        'auth0_github',
        'github|paul_atreides',
        true,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_3',
        'auth0',
        'auth0_username',
        'auth0|ellen_ripley',
        false,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_4',
        'auth0',
        'auth0_username',
        'auth0|ulfric_stormcloak',
        false,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_5',
        'auth0',
        'auth0_username',
        'auth0|chrisjen_avasarala',
        false,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_6',
        'auth0',
        'auth0_username',
        'auth0|amos_burton',
        false,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_7',
        'auth0',
        'auth0_username',
        'auth0|frodo_baggins',
        false,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_8',
        'auth0',
        'auth0_username',
        'auth0|rick_deckard',
        false,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_9',
        'auth0',
        'auth0_username',
        'auth0|tyrion_lannister',
        false,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'user_11',
        'auth0',
        'auth0_username',
        'auth0|bernard_marx',
        false,
        '2022-11-11 11:11:11.000000'
    );
-- Insert user_roles
INSERT INTO user_roles (
        role_name,
        permissions,
        created_at,
        updated_at
    )
VALUES (
        'OWNER',
        ARRAY ['read:endpoint'::permissions, 'write:endpoint'::permissions, 'delete:endpoint'::permissions, 'transfer:endpoint'::permissions],
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'ADMIN',
        ARRAY ['read:endpoint'::permissions, 'write:endpoint'::permissions],
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'MEMBER',
        ARRAY ['read:endpoint'::permissions],
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
-- Insert account_user_access
INSERT INTO account_user_access (
        account_id,
        user_id,
        role_name,
        accepted,
        created_at,
        updated_at
    )
VALUES (
        'account_1',
        'user_1',
        'OWNER',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_1',
        'user_2',
        'ADMIN',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_1',
        'user_8',
        'ADMIN',
        false,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_2',
        'user_3',
        'OWNER',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_2',
        'user_4',
        'MEMBER',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_2',
        'user_9',
        'MEMBER',
        false,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_2',
        'user_2',
        'MEMBER',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_3',
        'user_5',
        'OWNER',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_3',
        'user_6',
        'ADMIN',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_3',
        'user_7',
        'MEMBER',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_3',
        'user_10',
        'MEMBER',
        false,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_4',
        'user_4',
        'OWNER',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'account_5',
        'user_4',
        'OWNER',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO portal_applications (
        id,
        account_id,
        name,
        gigastake,
        staked,
        request_timeout,
        gigastake_redirect,
        first_date_surpassed,
        created_at,
        updated_at
    )
VALUES (
        'test_app_1',
        'account_1',
        'pokt_app_123',
        true,
        false,
        5000,
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2',
        'account_2',
        'pokt_app_456',
        false,
        true,
        10000,
        false,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_3',
        'account_3',
        'pokt_app_789',
        false,
        true,
        10000,
        false,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO stickiness_options (
        lb_id,
        duration,
        sticky_max,
        stickiness,
        origins
    )
VALUES (
        'test_app_1',
        60,
        300,
        true,
        '{ "chrome-extension://", "moz-extension://" }'
    ),
    (
        'test_app_2',
        30,
        600,
        true,
        '{ "https://example.com", "https://test.com" }'
    );
INSERT INTO portal_application_aats (
        id,
        application_id,
        address,
        public_key,
        private_key,
        client_public_key,
        signature,
        version
    )
VALUES (
        'test_protocol_app_1',
        'test_app_1',
        'test_34715cae753e67c75fbb340442e7de8e',
        'test_34715cae753e67c75fbb340442e7de8e',
        'test_11b8d394ca331d7c7a71ca1896d630f6',
        'test_89a3af6a587aec02cfade6f5000424c2',
        'test_1dc39a2e5a84a35bf030969a0b3231f7',
        '0.0.1'
    ),
    (
        'test_protocol_app_2',
        'test_app_2',
        'test_8237c72345f12d1b1a8b64a1a7f66fa4',
        'test_8237c72345f12d1b1a8b64a1a7f66fa4',
        'test_2e83c836a29b423a47d8e18c779fd422',
        'test_04c71d90a92f40416b6f1d7d8af17e02',
        'test_f48d33b30ddaf60a1e5bb50d2ba8da5a',
        '0.0.1'
    ),
    (
        'test_protocol_app_3',
        'test_app_3',
        'test_b5e07928fc80083c13ad0201b81bae9b',
        'test_f608500e4fe3e09014fe2411b4a560b5',
        'test_8663e187c19f3c6e27317eab4ed6d7d5',
        'test_328a9cf1b35085eeaa669aa858f6fba9',
        'test_c3cd8be16ba32e24dd49fdb0247fc9b8',
        '0.0.1'
    );
INSERT INTO portal_application_settings (
        application_id,
        secret_key,
        secret_key_required,
        monthly_relay_limit,
        environment,
        favorited_chain_ids,
        updated_at
    )
VALUES (
        'test_app_1',
        'test_40f482d91a5ef2300ebb4e2308c',
        true,
        2500000,
        'production',
        ARRAY ['0001', '0053'],
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2',
        'test_9c9e3b193cfba5348f93bb2f3e3fb794',
        false,
        1500000,
        'production',
        ARRAY ['0021', '0064'],
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_3',
        'test_9f48b13e2bc5fd31ab367841f11495c1',
        false,
        4500000,
        'production',
        ARRAY ['0001', '0034'],
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO portal_application_whitelists (
        application_id,
        type,
        value,
        chain_id,
        created_at
    )
VALUES (
        'test_app_1',
        'blockchains',
        '0053',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_1',
        'contracts',
        '0x1234567890abcdef',
        '0001',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_1',
        'methods',
        'GET',
        '0001',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_1',
        'origins',
        'https://test.com',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_1',
        'userAgents',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2',
        'blockchains',
        '0021',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2',
        'contracts',
        '0x0987654321abcdef',
        '0064',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2',
        'methods',
        'POST',
        '0064',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2',
        'origins',
        'https://example.com',
        NULL,
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2',
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
        updated_at
    )
VALUES (
        'test_app_1',
        true,
        'email',
        'test@test.com',
        'trigger123',
        ARRAY ['quarter'::notification_event,'threeQuarters'::notification_event,'full'::notification_event],
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_2',
        true,
        'email',
        'email@pokt.network',
        'trigger456',
        ARRAY ['full'::notification_event,'half'::notification_event],
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO chains (
        id,
        blockchain,
        description,
        enforce_result,
        path,
        ticker,
        chain_aliases,
        log_limit_blocks,
        request_timeout,
        active,
        created_at,
        updated_at
    )
VALUES (
        '0001',
        'mainnet',
        'Pocket Network Mainnet',
        'JSON',
        '/v1/query/height',
        'POKT',
        ARRAY ['mainnet'],
        0,
        0,
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0053',
        'optimism-mainnet',
        'Optimism Mainnet',
        'JSON',
        '',
        'OP',
        ARRAY ['optimism-mainnet'],
        100000,
        0,
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0021',
        'eth-mainnet',
        'Ethereum Mainnet',
        'JSON',
        '',
        'ETH',
        ARRAY ['eth-mainnet'],
        100000,
        0,
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0064',
        'sui-testnet',
        'Sui Testnet',
        'JSON',
        '',
        'SUI-TESTNET',
        ARRAY ['sui-testnet'],
        100000,
        60000,
        false,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0040',
        'harmony-0',
        'Harmony Shard 0',
        'JSON',
        '',
        'HMY',
        ARRAY ['harmony-0'],
        0,
        0,
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO chain_altruists (
        chain_id,
        url,
        auth,
        auth_type,
        created_at,
        updated_at
    )
VALUES (
        '0001',
        'https://altruist-0001.com:1234',
        'test_pocket:auth123456',
        'basic_auth',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0053',
        'https://altruist-0053.com:1234',
        'test_pocket:auth123456',
        'basic_auth',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0021',
        'https://altruist-0021.com:1234',
        'test_pocket:auth123456',
        'basic_auth',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0064',
        'https://altruist-0064.com:1234',
        'test_pocket:auth123456',
        'basic_auth',
        '2022-11-11 11:11:11.000000',
        '2023-01-23T23:52:00.176019Z'
    ),
    (
        '0040',
        'https://altruist-0040.com:1234',
        'test_pocket:auth123456',
        'basic_auth',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO chain_gigastake_redirects (
        chain_id,
        portal_application_id,
        alias,
        domain,
        created_at,
        updated_at
    )
VALUES (
        '0001',
        'test_app_1',
        'altruist-0001',
        'pokt-rpc.gateway.pokt.network',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0053',
        'test_app_2',
        'altruist-0053',
        'op-rpc.gateway.pokt.network',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0021',
        'test_app_3',
        'altruist-0021',
        'eth-rpc.gateway.pokt.network',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0040',
        'test_app_3',
        'altruist-0040',
        'hmy-rpc.gateway.pokt.network',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
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
        '0001',
        'sync',
        '{"id":1,"jsonrpc":"2.0","method":"query"}',
        'result.sync_info',
        1,
        NULL,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0053',
        'sync',
        '{"id":1,"jsonrpc":"2.0","method":"eth_blockNumber","params":[]}',
        'result',
        2,
        NULL,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0021',
        'sync',
        '{"id":1,"jsonrpc":"2.0","method":"eth_blockNumber","params":[]}',
        'result',
        5,
        NULL,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0021',
        'chain',
        '{"method":"eth_chainId","id":1,"jsonrpc":"2.0"}',
        'id',
        NULL,
        1,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0064',
        'sync',
        '{"id":1,"jsonrpc":"2.0","method":"sui_blockNumber","params":[]}',
        'result',
        7,
        NULL,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0040',
        'sync',
        '{"id":1,"jsonrpc":"2.0","method":"hmy_blockNumber","params":[]}',
        'result',
        8,
        NULL,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO global_blocked_contracts (blocked_address, created_at, updated_at)
VALUES (
        '0xtest_6789abcdef0123456789abcdef01234567',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0xtest_f0123456789abcdef0123456789abcdef01',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0xtest_cdef0123456789abcdef0123456789abcdef',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0xtest_56789abcdef0123456789abcdef01234567',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0xtest_789abcdef0123456789abcdef0123456789',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
