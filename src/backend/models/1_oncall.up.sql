create table oncall (
identifier varchar (255) not null,
team_name varchar (255) not null,
team_type varchar (255) not null,
primary key (identifier, team_type)
);
