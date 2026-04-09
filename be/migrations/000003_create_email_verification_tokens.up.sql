-- Migración: 000003_create_email_verification_tokens.up.sql
-- ¿Qué? Crea la tabla para tokens de verificación de email tras el registro.
-- ¿Para qué? Confirmar que el usuario tiene acceso real al email que proporcionó.
-- ¿Impacto? Sin esta tabla, no hay forma de verificar emails y el sistema
--            quedaría abierto a registros con emails falsos o inexistentes.

CREATE TABLE email_verification_tokens (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,
    used       BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Índice único en token: búsqueda rápida cuando el usuario hace click en el enlace de verificación.
CREATE UNIQUE INDEX idx_email_verification_tokens_token ON email_verification_tokens (token);

-- Índice en user_id: útil para buscar el token activo de un usuario.
CREATE INDEX idx_email_verification_tokens_user_id ON email_verification_tokens (user_id);
