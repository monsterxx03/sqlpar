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

## Features

- [x] Basic select statement
- [x] Where clause support AND, OR, NOT, operators: >, >=, <, <=, =, !=, <>
- [x] Limit clause
- [x] `show table` statement
- [] Order by clause
- [] Column As  syntax
- [] Support nest column
- [] Group by clause
- [] Aggregate function
- [] Column function
- [] Export result as csv
