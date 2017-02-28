drop table if exists Servers;
create table if not exists Servers(
	id serial primary key not null,
	title text,
	address text,
	environment text
);

insert into Servers (title, address, environment) values
	('MoB Dev', 'dev-wp-wise-fs-1.spindance.net', 'wp_dev'),
	('Products Dev', 'dev-wp-products-fs-1.spindance.net', 'wp_dev'),
	('Features Dev', 'dev-wp-features-fs-1.spindance.net', 'wp_dev');

drop table if exists Actions;
create table if not exists Actions(
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
	created_at timestamp,
	last_modified timestamp
);

alter table requests alter column created_at set default now();

create or replace function update_last_modified() returns trigger as $$
begin
    new.last_modified = now();
    return new;
end;
$$ language 'plpgsql';

create trigger update_last_modified
    before update
    on Requests
    for each row
    execute procedure update_last_modified();
