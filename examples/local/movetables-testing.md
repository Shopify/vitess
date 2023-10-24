```bash
source env.sh

# Bring up initial cluster and unsharded commerce keyspace
./101_initial_cluster.sh

# Bring up sharded customer keyspace shard -80
./201_newkeyspace_tablets.sh customer 2 -80

# Bring up sharded customer keyspace shard 80-
./201_newkeyspace_tablets.sh customer 3 80-

# Bring up sharded lookups keyspace shard -80
./201_newkeyspace_tablets.sh lookups 4 -80

# Bring up sharded lookups keyspace shard 80-
./201_newkeyspace_tablets.sh lookups 5 80-

# Apply vschemas: (commerce vschema was already applied in 101_initial_cluster)
vtctldclient ApplyVSchema --vschema-file vschema_customer_sharded.json customer
vtctldclient ApplyVSchema --vschema-file vschema_lookups_sharded.json lookups

# Apply schemas (commerce schema was already applied in 101_initial_cluster)
vtctldclient ApplySchema --sql-file create_customer_sharded.sql customer
vtctldclient ApplySchema --sql-file create_lookups_sharded.sql lookups

# Insert data
mysql < ../common/insert_commerce_data.sql

# Confirm lookups are being used
mysql --comments --execute "EXPLAIN format=vtexplain SELECT * FROM corder WHERE order_id = 1000"
+------+----------+-------+--------------------------------------------------------------------------------+
| #    | keyspace | shard | query                                                                          |
+------+----------+-------+--------------------------------------------------------------------------------+
|    0 | lookups  | -80   | select order_id, keyspace_id from corder_keyspace_idx where order_id in (1000) |
|    1 | customer | -80   | select order_id, customer_id, sku, price from corder where order_id = 1000     |
+------+----------+-------+--------------------------------------------------------------------------------+
2 rows in set (0.00 sec)

# Bring up sharded lookups2 keyspace shard -80
./201_newkeyspace_tablets.sh lookups2 6 -80

# Bring up sharded lookups2 keyspace shard 80-
./201_newkeyspace_tablets.sh lookups2 7 80-

# Apply original lookups vschema to lookups2
vtctldclient ApplyVSchema --vschema-file vschema_lookups_sharded.json lookups2

# Start moving lookups to lookups2
vtctlclient MoveTables -- --source lookups --tables='corder_keyspace_idx' Create lookups2.lookups_to_lookups2

# Ensure copy completed and vreplication is running
vtctlclient MoveTables -- progress lookups2.lookups_to_lookups2

# Start vdiff
vtctlclient VDiff -- --v2 lookups2.lookups_to_lookups2 create

# Confirm success
vtctlclient VDiff -- --v2 lookups2.lookups_to_lookups2 show last

# Switch traffic
vtctlclient MoveTables -- SwitchTraffic lookups2.lookups_to_lookups2
 
# Confirm queries are using the new lookups2 keyspace
mysql --comments --execute "EXPLAIN format=vtexplain SELECT * FROM corder WHERE order_id = 1000"
+------+----------+-------+--------------------------------------------------------------------------------+
| #    | keyspace | shard | query                                                                          |
+------+----------+-------+--------------------------------------------------------------------------------+
|    0 | lookups2 | -80   | select order_id, keyspace_id from corder_keyspace_idx where order_id in (1000) |
|    1 | customer | -80   | select order_id, customer_id, sku, price from corder where order_id = 1000     |
+------+----------+-------+--------------------------------------------------------------------------------+
2 rows in set (0.00 sec)

# Confirm inserts will go to the new lookups2 keyspace
mysql --comments --execute "EXPLAIN /*vt+ EXECUTE_DML_QUERIES */  format=vtexplain insert into corder(customer_id, sku, price) values(5, 'SKU-1002', 30)"
+------+----------+-------+----------------------------------------------------------------------------------------+
| #    | keyspace | shard | query                                                                                  |
+------+----------+-------+----------------------------------------------------------------------------------------+
|    0 | lookups2 | -80   | begin                                                                                  |
|    0 | lookups2 | -80   | insert into corder_keyspace_idx(order_id, keyspace_id) values (1009, 'p�<�
                                                                                                      �z')          |
|    1 | customer | -80   | begin                                                                                  |
|    1 | customer | -80   | insert into corder(customer_id, sku, price, order_id) values (5, 'SKU-1002', 30, 1009) |
+------+----------+-------+----------------------------------------------------------------------------------------+

# Update the customer vschema to use `lookups2`
vtctlclient ApplyVSchema -- --vschema "$(cat vschema_customer_sharded.json | sed s/lookups\.corder_keyspace_idx/lookups2\.corder_keyspace_idx/)" customer

# Complete the movetables
vtctlclient MoveTables -- Complete lookups2.lookups_to_lookups2

# Confirm you can insert and query records
mysql --comments

mysql> EXPLAIN format=vtexplain SELECT * FROM corder WHERE order_id = (SELECT MAX(order_id) FROM corder);
+------+----------+-------+--------------------------------------------------------------------------------+
| #    | keyspace | shard | query                                                                          |
+------+----------+-------+--------------------------------------------------------------------------------+
|    0 | customer | -80   | select max(order_id) from corder                                               |
|    1 | customer | 80-   | select max(order_id) from corder                                               |
|    2 | lookups2 | 80-   | select order_id, keyspace_id from corder_keyspace_idx where order_id in (1101) |
|    3 | customer | -80   | select order_id, customer_id, sku, price from corder where order_id = 1101     |
+------+----------+-------+--------------------------------------------------------------------------------+
4 rows in set (0.00 sec)

mysql> EXPLAIN /*vt+ EXECUTE_DML_QUERIES */  format=vtexplain insert into corder(customer_id, sku, price) values(5, 'SKU-1002', 30);
+------+----------+-------+----------------------------------------------------------------------------------------+
| #    | keyspace | shard | query                                                                                  |
+------+----------+-------+----------------------------------------------------------------------------------------+
|    0 | lookups2 | -80   | begin                                                                                  |
|    0 | lookups2 | -80   | insert into corder_keyspace_idx(order_id, keyspace_id) values (1102, 'p�<�
                                                                                                      �z')          |
|    1 | customer | -80   | begin                                                                                  |
|    1 | customer | -80   | insert into corder(customer_id, sku, price, order_id) values (5, 'SKU-1002', 30, 1102) |
+------+----------+-------+----------------------------------------------------------------------------------------+
4 rows in set (0.01 sec)

mysql> SELECT * FROM corder WHERE order_id = (SELECT MAX(order_id) FROM corder);
+----------+-------------+--------------------+-------+
| order_id | customer_id | sku                | price |
+----------+-------------+--------------------+-------+
|     1102 |           5 | 0x534B552D31303032 |    30 |
+----------+-------------+--------------------+-------+
1 row in set (0.00 sec)

# Confirm table was removed from original lookups keyspace
mysql> show tables from lookups;
Empty set (0.00 sec)

# Confirm table removed from original lookups vschema
vtctlclient GetVschema lookups
{
  "sharded": true,
  "vindexes": {
    "hash": {
      "type": "hash"
    }
  }
}
```
