create table if not exists Requests(
	id serial primary key not null,
	requester text not null,
	action text,
	server text,
	responder text,
	approved bool,
	success bool,
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

ALTER TABLE Requests ADD COLUMN response_url text;
