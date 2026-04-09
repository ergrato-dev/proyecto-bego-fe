-- Migración: 000002_create_password_reset_tokens.down.sql
-- ¿Qué? Revierte la creación de la tabla password_reset_tokens.
-- ¿Para qué? Deshacer la migración en caso de error durante el desarrollo.
-- ¿Impacto? ELIMINA TODOS LOS tokens de recuperación existentes.

DROP TABLE IF EXISTS password_reset_tokens;
