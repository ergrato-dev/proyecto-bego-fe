-- Migración: 000003_create_email_verification_tokens.down.sql
-- ¿Qué? Elimina la tabla de tokens de verificación de email.
-- ¿Para qué? Revertir la migración 000003 de forma limpia sin dejar tablas huérfanas.
-- ¿Impacto? Al hacer rollback, los tokens de verificación existentes se pierden.

DROP TABLE IF EXISTS email_verification_tokens;
