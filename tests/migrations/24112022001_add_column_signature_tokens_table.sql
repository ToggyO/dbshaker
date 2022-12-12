-- +dbshaker UpStart

ALTER TABLE tokens
    ADD COLUMN signature VARCHAR(100);

-- +dbshaker UpEnd

-- +dbshaker DownStart

ALTER TABLE tokens
    DROP COLUMN signature;

-- +dbshaker DownEnd