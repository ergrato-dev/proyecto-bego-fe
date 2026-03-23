# HU-001 — Registro de Cuenta

<!--
  ¿Qué? Historia de usuario para el registro de nuevas cuentas.
  ¿Para qué? Definir el comportamiento esperado del sistema durante el registro.
  ¿Impacto? Sin esta HU, el equipo no tendría criterios claros de aceptación
             para el flujo de registro, generando implementaciones inconsistentes.
-->

## Identificación

| Campo         | Valor                                                               |
| ------------- | ------------------------------------------------------------------- |
| **ID**        | HU-001                                                              |
| **Módulo**    | Autenticación                                                       |
| **Prioridad** | Alta                                                                |
| **Estado**    | Pendiente                                                           |
| **Sprint**    | Sprint 1                                                            |
| **RFs**       | RF-001                                                              |
| **RNFs**      | RNF-001 (Seguridad), RNF-002 (Rendimiento), RNF-004 (Accesibilidad) |

---

## Historia

**Como** visitante del sistema (usuario no autenticado),
**quiero** poder crear una cuenta con mi email y contraseña,
**para** acceder a las funcionalidades del sistema NN Auth.

---

## Criterios de Aceptación

### CA-001.1 — Registro exitoso

**Dado** que soy un visitante en la página de registro
**Y** completo el formulario con email único, nombre completo y contraseña válida
**Cuando** hago clic en el botón "Crear cuenta"
**Entonces** el sistema crea mi cuenta con `is_email_verified = false`
**Y** envía un email de verificación a la dirección proporcionada
**Y** muestra un mensaje de éxito: "Cuenta creada. Revisa tu email para verificar tu cuenta."
**Y** no inicia sesión automáticamente (requiere verificación de email primero)

### CA-001.2 — Email ya registrado

**Dado** que existe un usuario con el email "maria@ejemplo.com"
**Cuando** intento registrarme con el mismo email
**Entonces** el sistema muestra el mensaje de error "Ya existe una cuenta con este email"
**Y** no crea un nuevo usuario
**Y** el formulario permanece con los datos ingresados

### CA-001.3 — Contraseña débil

**Dado** que estoy en el formulario de registro
**Cuando** ingreso una contraseña que no cumple los requisitos mínimos (< 8 chars, sin mayúscula, sin número)
**Entonces** el sistema muestra el error de validación específico al campo de contraseña
**Y** el formulario no se envía
**Y** el foco permanece en el campo de contraseña

### CA-001.4 — Campos obligatorios vacíos

**Dado** que el formulario de registro tiene campos vacíos
**Cuando** intento enviar el formulario
**Entonces** el sistema muestra un mensaje de error en cada campo vacío obligatorio
**Y** no realiza ninguna petición al servidor

### CA-001.5 — Email con formato inválido

**Dado** que ingreso un email con formato inválido (ej: "no-es-un-email")
**Cuando** intento enviar el formulario
**Entonces** el sistema muestra "Ingresa un email válido" en el campo de email

### CA-001.6 — Seguridad: contraseña no visible en ninguna respuesta

**Dado** que me registro exitosamente
**Cuando** el servidor responde
**Entonces** la respuesta JSON NO contiene el campo `password` ni `hashed_password`

### CA-001.7 — Accesibilidad del formulario

**Dado** que un usuario navega el formulario solo con teclado
**Entonces** puede completar y enviar el formulario usando Tab, Enter y Shift+Tab
**Y** los mensajes de error son anunciados por lectores de pantalla (`role="alert"`)

---

## Notas Técnicas

- El endpoint es `POST /api/v1/auth/register`
- La contraseña se hashea con bcrypt (golang.org/x/crypto/bcrypt) antes de almacenarse
- El email de verificación usa un token UUID con expiración de 24 horas
- El token se almacena en la tabla `email_verification_tokens`
- El handler Go usa `go-playground/validator` para validar el DTO de request
- Rate limiting: 5 peticiones / 15 minutos por IP
