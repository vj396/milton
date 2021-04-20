create table interrupts (
id int unsigned not null auto_increment,
item varchar (255) not null,
submitted_by varchar (255) not null,
submitted_at bigint unsigned not null,
channel_id varchar (255) not null,
primary key (id),
foreign key (channel_id) references membership(id) on delete cascade
);