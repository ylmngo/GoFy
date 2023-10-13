CREATE TABLE IF NOT EXISTS files (
    file_id bytea PRIMARY KEY, 
    user_id bigint REFERENCES users ON DELETE CASCADE,  
    filename text NOT NULL, 
    metadata text NOT NULL,  
    uploaded_at timestamp(0) with time zone NOT NULL DEFAULT NOW()  
);                                                                      