// Archivo: user.go
// Descripción: Modelo GORM que representa la tabla `users` en PostgreSQL.
// ¿Para qué? Proveer la estructura Go que GORM usa para queries de usuarios.
// ¿Impacto? Cada campo mal tipado o mal tagueado genera queries SQL incorrectas
//            o falla silenciosamente — revisar siempre los tags de GORM.

package models

import "time"

// User representa un usuario registrado en el sistema.
// ¿Para qué? Mapear la tabla `users` de la base de datos a un struct de Go.
// ¿Impacto? Este struct es el centro del sistema de auth — todas las operaciones
//
//	de registro, login y cambio de contraseña usan este modelo.
type User struct {
	// ¿Qué? UUID generado por PostgreSQL como clave primaria.
	// ¿Para qué? Los UUIDs son globalmente únicos y no revelan el volumen de registros
	//            (a diferencia de IDs secuenciales como 1, 2, 3...).
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	// ¿Qué? Email único del usuario — también sirve como identificador de login.
	// ¿Para qué? El índice (uniqueIndex) garantiza que no haya duplicados en BD
	//            y acelera la búsqueda por email en el login.
	Email string `gorm:"type:varchar(255);uniqueIndex;not null"`

	// ¿Qué? Nombre completo del usuario para mostrar en la interfaz.
	FullName string `gorm:"type:varchar(255);not null"`

	// ¿Qué? Hash bcrypt de la contraseña — NUNCA la contraseña en texto plano.
	// ¿Para qué? Almacenar la contraseña de forma que sea imposible revertirla
	//            a texto plano, incluso si la base de datos es comprometida.
	// ¿Impacto? Si se almacena el texto plano, una filtración expone todas las
	//            contraseñas de los usuarios directamente.
	HashedPassword string `gorm:"type:varchar(255);not null"`

	// ¿Qué? Indica si la cuenta está activa.
	// ¿Para qué? Permite desactivar cuentas sin eliminarlas (soft disable).
	IsActive bool `gorm:"default:true;not null"`

	// ¿Qué? Indica si el usuario verificó su email tras el registro.
	// ¿Para qué? Bloquear el login hasta que el email sea verificado,
	//            evitando registros con emails falsos o mal escritos.
	IsEmailVerified bool `gorm:"default:false;not null"`

	// ¿Qué? Timestamps gestionados automáticamente por GORM.
	// ¿Para qué? Registro de auditoría — saber cuándo se creó/modificó cada usuario.
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`
}
