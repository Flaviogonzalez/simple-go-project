-- mismo id sin autoincrement, y encima en la query tengo que añadir el id manualmente

CREATE TYPE logtype AS ENUM ('debug', 'info', 'warning', 'error', 'fatal');

CREATE TABLE log (
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL,
    thread_identifier TEXT,
    requestid TEXT,
    logtype logtype NOT NULL,
    userid TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
); 