CREATE TABLE IF NOT EXISTS items (
    id integer not null primary key,
    title text not null,
    body text,
    priority integer not null,
    done bool,
    date_complete integer
);
