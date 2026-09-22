CREATE TABLE link_clicks (
    id BIGSERIAL PRIMARY KEY,
    link_id BIGINT NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    clicked_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX link_clicks_link_id_clicked_at_idx ON link_clicks (link_id, clicked_at);
