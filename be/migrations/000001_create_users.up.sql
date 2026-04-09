-- Migración: 000001_create_users.up.sql
-- ¿Qué? Crea la tabla principal de usuarios del sistema.
-- ¿Para qué? Almacenar las cuentas de usuario con sus credenciales hasheadas.
-- ¿Impacto? Sin esta tabla, ningún flujo de autenticación puede funcionar.

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id               UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    email            VARCHAR(255) NOT NULL,
    full_name        VARCHAR(255) NOT NULL,
    hashed_password  VARCHAR(255) NOT NULL,
    is_active        BOOLEAN      NOT NULL DEFAULT TRUE,
    is_email_verified BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Índice único en email: acelera el login (búsqueda por email) y garantiza unicidad.
-- Sin este índice, un SELECT WHERE email = ? haría full table scan.
CREATE UNIQUE INDEX idx_users_email ON users (email);
