drop table if exists onlineddl_test;
create table onlineddl_test (
  id bigint not null auto_increment,
  i int not null,
  ts timestamp(6),
  primary key(id),
  constraint `check_1` CHECK ((`i` >= 0)),
  constraint `check_2` CHECK ((`i` <> 10)),
  constraint `check_3` CHECK ((`i` >= 0))
) ;

drop event if exists onlineddl_test;
delimiter ;;
create event onlineddl_test
  on schedule every 1 second
  starts current_timestamp
  ends current_timestamp + interval 60 second
  on completion not preserve
  enable
  do
begin
  insert into onlineddl_test values (null, 11, now(6));
  insert into onlineddl_test values (null, 13, now(6));
  insert into onlineddl_test values (null, 17, now(6));
  insert into onlineddl_test values (null, 19, now(6));
end ;;
