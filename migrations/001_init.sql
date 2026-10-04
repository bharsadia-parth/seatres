CREATE TABLE IF NOT EXISTS shows (
  id             UUID PRIMARY KEY,
  name           TEXT   NOT NULL,
  price_paise    BIGINT NOT NULL,
  per_user_limit INT    NOT NULL DEFAULT 4
);

CREATE TABLE IF NOT EXISTS reservations (
  id           UUID PRIMARY KEY,
  show_id      UUID   NOT NULL REFERENCES shows(id),
  user_id      TEXT   NOT NULL,
  seats        TEXT[] NOT NULL,
  amount_paise BIGINT NOT NULL,
  status       TEXT   NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS seats (
  show_id        UUID NOT NULL REFERENCES shows(id),
  label          TEXT NOT NULL,
  status         TEXT NOT NULL DEFAULT 'available',
  user_id        TEXT,
  reservation_id UUID,
  PRIMARY KEY (show_id, label),
  CHECK (status IN ('available','held','confirmed'))
);

CREATE INDEX IF NOT EXISTS seats_user_idx ON seats(show_id, user_id);

CREATE TABLE IF NOT EXISTS idempotency_keys (
  user_id        TEXT NOT NULL,
  key            TEXT NOT NULL,
  request_hash   TEXT NOT NULL,
  reservation_id UUID,
  PRIMARY KEY (user_id, key)
);