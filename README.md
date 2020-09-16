# sqlpar

Use sql to query parquet file

        ./sqlpar -file test.parquet

        >> select user_id, age from test_schema where age > 1 limit 2
        user_id  age
        abc      12
        def      33

        ./sqlpar -file test.parquet -sql 'select user_id, age from test_schema where age > 1 limit 2'
        user_id  age
        abc      12
        def      33
