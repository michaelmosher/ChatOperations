drop table if exists Servers;
create table Servers(
	id serial primary key not null,
	title text,
	address text,
	environment text
);

insert into Servers (title, address, environment) values
	('MoB Dev', 'dev-wp-wise-fs-1.spindance.net', 'wp_dev'),
	('Products Dev', 'dev-wp-products-fs-1.spindance.net', 'wp_dev');

drop table if exists Actions;
create table Actions(
	id serial primary key not null,
	title text,
	command text
);

insert into Actions (title, command) values
	('Deploy', 'sudo chef-client'),
	('Config Loader', 'sudo -E -u onewise bundle exec bin/rake configuration:load');

create table if not exists Requests(
	id serial primary key not null,
	requester text not null,
	actionId int4 references Actions(id),
	serverId int4 references Servers(id),
	responder text,
	approved bool,
	success bool,
	response_url text,
	created_at datetime default CURRENT_TIMESTAMP,
	last_modified datetime on update CURRENT_TIMESTAMP
);
