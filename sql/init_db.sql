create table IF NOT EXISTS links (short varchar(8), long varchar(1000));
create index concurrently links_short_idx on links using btree (short text_pattern_ops);

create table IF NOT EXISTS stats (link varchar(8), ip inet);
create index concurrently stats_link_ip_idx on stats using btree (link text_pattern_ops, ip inet_ops);