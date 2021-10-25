create table IF NOT EXISTS links (short varchar(8), long varchar(1000));
create index concurrently links_short_idx on links using btree (short text_pattern_ops);