SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'cash_flow' AND pid <> pg_backend_pid();
DROP DATABASE "cash_flow";

CREATE DATABASE "cash_flow" WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';

ALTER DATABASE "cash_flow" OWNER TO "postgres";

\connect "cash_flow"

CREATE TABLE public.users (
	id bigserial NOT NULL,
	name character varying(255) NOT NULL,
	email character varying(255) NOT NULL,
	password_hash character varying(255) NOT NULL,
	password_token character(64),
	activation_token character(64),
	activated_at timestamp without time zone,
	time_zone character varying(20) NOT NULL,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone,
	deleted_at timestamp without time zone,
	PRIMARY KEY (id),
	UNIQUE (email),
	UNIQUE (password_token),
	UNIQUE (activation_token)
);

CREATE INDEX ON public.users (name);
CREATE INDEX ON public.users (activated_at);
CREATE INDEX ON public.users (created_at);
CREATE INDEX ON public.users (updated_at);
CREATE INDEX ON public.users (deleted_at);

CREATE TABLE public.journals (
	id bigserial NOT NULL,
	user_id bigint NOT NULL,
	item character varying(255) NOT NULL,
	datetime timestamp without time zone NOT NULL,
	amount integer NOT NULL,
	category smallint NOT NULL,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone,
	deleted_at timestamp without time zone,
	PRIMARY KEY (id),
	FOREIGN KEY (user_id) REFERENCES public.users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX ON public.journals (item);
CREATE INDEX ON public.journals (datetime);
CREATE INDEX ON public.journals (category);
CREATE INDEX ON public.journals (created_at);
CREATE INDEX ON public.journals (updated_at);
CREATE INDEX ON public.journals (deleted_at);