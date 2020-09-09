# sqlpar
Use sql to query parquet file

        ./sqlpar test.parquet

        >> select user_id, age from test_schema limit 2
        user_id  age
        abc      12
        def      33
