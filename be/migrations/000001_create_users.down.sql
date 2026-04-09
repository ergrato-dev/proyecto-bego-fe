-- Migración: 000001_create_users.down.sql
-- ¿Qué? Revierte la creación de la tabla users.
-- ¿Para qué? Permitir deshacer la migración en caso de error o durante el desarrollo.
-- ¿Impacto? ELIMINA TODOS LOS DATOS de usuarios — solo usar en desarrollo.

DROP TABLE IF EXISTS users;
