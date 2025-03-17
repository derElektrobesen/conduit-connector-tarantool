return {
    sharding = {
        ['cbf06940-0790-498b-948d-042b62cf3d29'] = { -- replicaset #1
            master = 'auto',
            replicas = {
                ['8a274925-a26d-47fc-9e1b-af88ce939412'] = {
                    uri = 'storage:storage@storage_1_a:3301',
                    name = 'storage_1_a',
                },
                ['3de2e3e1-9ebe-4d0d-abb1-26d301b84633'] = {
                    uri = 'storage:storage@storage_1_b:3302',
                    name = 'storage_1_b'
                },
                ['3de2e3e1-9ebe-4d0d-abb1-26d301b84635'] = {
                    uri = 'storage:storage@storage_1_c:3303',
                    name = 'storage_1_c'
                }
            },
        }, -- replicaset #1
        ['ac522f65-aa94-4134-9f64-51ee384f1a54'] = { -- replicaset #2
            master = 'auto',
            replicas = {
                ['1e02ae8a-afc0-4e91-ba34-843a356b8ed7'] = {
                    uri = 'storage:storage@storage_2_a:3311',
                    name = 'storage_2_a',
                },
                ['001688c3-66f8-4a31-8e19-036c17d489c2'] = {
                    uri = 'storage:storage@storage_2_b:3312',
                    name = 'storage_2_b'
                },
                ['c23516d5-22de-4ef4-8918-73d52e7661e2'] = {
                    uri = 'storage:storage@storage_2_c:3313',
                    name = 'storage_2_c'
                }
            },
        }, -- replicaset #2
    }, -- sharding
    replication_connect_quorum = 0,
    election_mode = "candidate"
    bucket_count = 10000,
}
