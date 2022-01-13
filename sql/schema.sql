CREATE TABLE devices (
     ieee_address varchar(24) not null,
     date_created timestamp with time zone not null,
     date_modified timestamp with time zone not null,
     date_code varchar(24) null,
     friendly_name varchar(255) not null,
     area_id bigint references areas(id),
     manufacturer text null,
     model_id text null,
     last_seen timestamp with time zone null,
     type varchar(64) not null,
     battery numeric null,
     active bool not null
);

ALTER TABLE devices ADD PRIMARY KEY (ieee_address);

CREATE TABLE areas (
   id serial not null,
   uuid uuid not null,
   date_created timestamp with time zone not null,
   name varchar(255) not null
);


ALTER TABLE areas ADD PRIMARY KEY (id);
ALTER TABLE areas ADD UNIQUE (uuid);


CREATE TABLE temperature_reports (
     device_id varchar(24) not null references devices(ieee_address),
     area_id bigint not null references areas(id),
     date timestamp with time zone not null,
     value numeric null
);

CREATE TABLE pressure_reports (
      device_id varchar(24) not null references devices(ieee_address),
      area_id bigint not null references areas(id),
      date timestamp with time zone not null,
      value numeric null
);


CREATE TABLE humidity_reports (
      device_id varchar(24) not null references devices(ieee_address),
      area_id bigint not null references areas(id),
      date timestamp with time zone not null,
      value numeric null
);

CREATE TABLE illuminance_reports (
     device_id varchar(24) not null references devices(ieee_address),
     area_id bigint not null references areas(id),
     date timestamp with time zone not null,
     value numeric null,
     value_lux numeric null
);

CREATE TABLE groups (
    id numeric not null,
    date_created timestamp with time zone not null,
    date_modified timestamp with time zone not null,
    friendly_name varchar(255) not null,
    active bool not null
);

ALTER TABLE groups ADD PRIMARY KEY (id);

CREATE TABLE groups_devices (
    group_id numeric not null,
    ieee_address varchar(24) not null
);

ALTER TABLE groups_devices ADD PRIMARY KEY (group_id, ieee_address);