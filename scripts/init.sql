CREATE SCHEMA IF NOT EXISTS users;

create table users.users(
	id uuid primary key,
	username text unique not null,
	password_hash text,
	name text not null,
	surname text not null,
	score integer,
	created_at timestamp default current_timestamp
);

CREATE SCHEMA IF NOT EXISTS posts;



CREATE SCHEMA IF NOT EXISTS comments;
CREATE SCHEMA IF NOT EXISTS likes;

GRANT ALL PRIVILEGES ON SCHEMA users TO postgres;
GRANT ALL PRIVILEGES ON SCHEMA posts TO postgres;
GRANT ALL PRIVILEGES ON SCHEMA comments TO postgres;
GRANT ALL PRIVILEGES ON SCHEMA likes TO postgres;