create table corder_keyspace_idx(
    order_id bigint,
    keyspace_id varbinary(10),
    primary key(order_id)
) ENGINE=InnoDB;

