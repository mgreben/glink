CREATE TABLE IF NOT EXISTS link_clicks (
    link_id Int64,
    clicked_at DateTime64(3, 'UTC')
) ENGINE = MergeTree
ORDER BY (link_id, clicked_at);
