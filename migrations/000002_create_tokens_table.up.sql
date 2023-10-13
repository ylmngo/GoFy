CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY, 
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,  
    expiry timestamp(0) with time zone not null, 
    scope text not null 
);