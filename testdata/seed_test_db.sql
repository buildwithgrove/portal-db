INSERT INTO pay_plans (plan_type, daily_limit)
VALUES ('FREETIER_V0', 250000),
    ('PAY_AS_YOU_GO_V0', 0),
    ('ENTERPRISE', 0),
    ('TEST_PLAN_V0', 100),
    ('TEST_PLAN_10K', 10000),
    ('TEST_PLAN_90K', 90000);
INSERT INTO user_roles (name, permissions)
VALUES ('ADMIN', '{ "read:endpoint", "write:endpoint" }'),
    ('OWNER', '{ "read:endpoint", "write:endpoint" }'),
    ('MEMBER', '{ "read:endpoint" }');
INSERT INTO applications (
        application_id,
        name,
        status,
        url,
        user_id,
        dummy,
        created_at,
        updated_at
    )
VALUES (
        'test_app_47hfnths73j2se',
        'pokt_app_123',
        'IN_SERVICE',
        'https://test.app123.io',
        'test_user_1dbffbdfeeb225',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_app_5hdf7sh23jd828',
        'pokt_app_456',
        'IN_SERVICE',
        'https://test.app456.io',
        'test_user_04228205bd261a',
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO app_limits (
        application_id,
        pay_plan,
        custom_limit
    )
VALUES (
        'test_app_47hfnths73j2se',
        'FREETIER_V0',
        null
    ),
    (
        'test_app_5hdf7sh23jd828',
        'ENTERPRISE',
        2000000
    );
INSERT INTO gateway_aat (
        application_id,
        address,
        public_key,
        signature,
        client_public_key,
        private_key
    )
VALUES (
        'test_app_47hfnths73j2se',
        'test_34715cae753e67c75fbb340442e7de8e',
        'test_11b8d394ca331d7c7a71ca1896d630f6',
        'test_89a3af6a587aec02cfade6f5000424c2',
        'test_1dc39a2e5a84a35bf030969a0b3231f7',
        'test_d2ce53f115f4ecb2208e9188800a85cf'
    ),
    (
        'test_app_5hdf7sh23jd828',
        'test_558c0225c7019e14ccf2e7379ad3eb50',
        'test_96c981db344ab6920b7e87853838e285',
        'test_1272a8ab4cbbf636f09bf4fa5395b885',
        'test_d709871777b89ed3051190f229ea3f01',
        'test_53e50765d8bc1fb41b3b0065dd8094de'
    );
INSERT INTO gateway_settings (
        application_id,
        secret_key,
        secret_key_required
    )
VALUES (
        'test_app_47hfnths73j2se',
        'test_40f482d91a5ef2300ebb4e2308c',
        true
    ),
    (
        'test_app_5hdf7sh23jd828',
        'test_90210ac4bdd3423e24877d1ff92',
        false
    );
INSERT INTO notification_settings (
        application_id,
        signed_up,
        on_quarter,
        on_half,
        on_three_quarters,
        on_full
    )
VALUES (
        'test_app_47hfnths73j2se',
        true,
        false,
        false,
        true,
        true
    ),
    (
        'test_app_5hdf7sh23jd828',
        true,
        false,
        false,
        true,
        true
    );
INSERT INTO loadbalancers (
        lb_id,
        user_id,
        name,
        request_timeout,
        gigastake,
        gigastake_redirect,
        created_at,
        updated_at
    )
VALUES (
        'test_lb_34987u329rfn23f',
        'test_user_1dbffbdfeeb225',
        'pokt_app_123',
        5000,
        true,
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_lb_3890ru23jfi32fj',
        'test_user_04228205bd261a',
        'pokt_app_456',
        5000,
        true,
        true,
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        'test_lb_34gg4g43g34g5hh',
        'test_user_redirect233344',
        'test_lb_redirect',
        5000,
        false,
        false,
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
        'test_lb_34987u329rfn23f',
        60,
        300,
        true,
        '{ "chrome-extension://", "moz-extension://" }'
    ),
    (
        'test_lb_3890ru23jfi32fj',
        40,
        400,
        true,
        '{ "chrome-extension://" }'
    ),
    (
        'test_lb_34gg4g43g34g5hh',
        20,
        600,
        false,
        '{ "test-extension://", "test-extension2://" }'
    );
INSERT INTO user_access (
        lb_id,
        role_name,
        user_id,
        email,
        accepted
    )
VALUES (
        'test_lb_34987u329rfn23f',
        'OWNER',
        'test_user_1dbffbdfeeb225',
        'owner1@test.com',
        true
    ),
    (
        'test_lb_34987u329rfn23f',
        'ADMIN',
        'test_user_admin1234',
        'admin1@test.com',
        true
    ),
    (
        'test_lb_34987u329rfn23f',
        'MEMBER',
        'test_user_member1234',
        'member1@test.com',
        true
    ),
    (
        'test_lb_3890ru23jfi32fj',
        'OWNER',
        'test_user_04228205bd261a',
        'owner2@test.com',
        true
    ),
    (
        'test_lb_3890ru23jfi32fj',
        'ADMIN',
        'test_user_admin5678',
        'admin2@test.com',
        true
    ),
    (
        'test_lb_34gg4g43g34g5hh',
        'OWNER',
        'test_user_redirect233344',
        'owner3@test.com',
        true
    ),
    (
        'test_lb_34gg4g43g34g5hh',
        'MEMBER',
        'test_user_member5678',
        'member2@test.com',
        true
    );
INSERT INTO lb_apps (lb_id, app_id)
VALUES (
        'test_lb_34987u329rfn23f',
        'test_app_47hfnths73j2se'
    ),
    (
        'test_lb_3890ru23jfi32fj',
        'test_app_5hdf7sh23jd828'
    );
INSERT INTO blockchains (
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
        '0001',
        true,
        'https://test:test_93uhfniu23f8@shared-test2.nodes.pokt.network:12345',
        'pokt-mainnet',
        '{ "pokt-mainnet" }',
        null,
        null,
        'POKT Network Mainnet',
        'JSON',
        100000,
        'POKT-mainnet',
        '',
        null,
        'POKT',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0021',
        true,
        'https://test:test_u32fh239hf@shared-test2.nodes.eth.network:12345',
        'eth-mainnet',
        '{ "eth-mainnet" }',
        '1',
        '{\"method\":\"eth_chainId\",\"id\":1,\"jsonrpc\":\"2.0\"}',
        'Ethereum Mainnet',
        'JSON',
        100000,
        'ETH-1',
        '',
        null,
        'ETH',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO redirects (
        blockchain_id,
        alias,
        loadbalancer,
        domain,
        created_at,
        updated_at
    )
VALUES (
        '0001',
        'test-mainnet',
        'test_lb_34gg4g43g34g5hh',
        'test-rpc1.testnet.pokt.network',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0001',
        'test-mainnet',
        'test_lb_34gg4g43g34g5hh',
        'test-rpc2.testnet.pokt.network',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    ),
    (
        '0021',
        'eth-mainnet',
        'test_lb_34gg4g43g34g5hh',
        'test-rpc.testnet.eth.network',
        '2022-11-11 11:11:11.000000',
        '2022-11-11 11:11:11.000000'
    );
INSERT INTO sync_check_options (
        blockchain_id,
        synccheck,
        allowance,
        body,
        path,
        result_key
    )
VALUES (
        '0001',
        null,
        1,
        '{}',
        '/v1/query/height',
        'height'
    ),
    (
        '0021',
        null,
        5,
        '{\"method\":\"eth_blockNumber\",\"id\":1,\"jsonrpc\":\"2.0\"}',
        null,
        'result'
    );
