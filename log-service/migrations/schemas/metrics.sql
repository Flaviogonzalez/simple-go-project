
-- shit code, id sin autoincrement, a total waste of time

CREATE TABLE metrics (
    id SERIAL PRIMARY KEY,
    histograms JSONB,
    counts INTEGER[],
    values DOUBLE PRECISION[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    requestid TEXT
);


CREATE TABLE request (
    id SERIAL PRIMARY KEY,
    requestid TEXT,
    code INTEGER,
    size INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
); 