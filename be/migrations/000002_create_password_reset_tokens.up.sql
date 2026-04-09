-- Migración: 000002_create_password_reset_tokens.up.sql
-- ¿Qué? Crea la tabla para tokens de recuperación de contraseña.
-- ¿Para qué? Almacenar tokens temporales enviados por email para restablecer contraseñas.
-- ¿Impacto? Sin esta tabla, el flujo "olvidé mi contraseña" no puede funcionar.

CREATE TABLE password_reset_tokens (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,
    used       BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Índice único en token: búsqueda rápida cuando el usuario envía el token del enlace.
CREATE UNIQUE INDEX idx_password_reset_tokens_token ON password_reset_tokens (token);

-- Índice en user_id: útil para invalidar todos los tokens de un usuario si cambia contraseña.
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens (user_id);
