CREATE TABLE channels_logs (
                               id text PRIMARY KEY,
                               channel_id text NOT NULL,
                               date timestamp without time zone NOT NULL,
                               error text
);
