package back_go

import (
	_ "github.com/supchat-lmrt/back-go/internal/models" // Importez pour que Swagger trouve les modèles
)

// @title        SupChat API
// @version      1.0
// @description  API pour le service de messagerie SupChat

// @contact.name  Support API SupChat
// @contact.email supinfo@supinfo.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name session_token

// @host      localhost:3000
// @BasePath
func SwaggerInfo() {
	// Cette fonction ne fait rien, elle sert uniquement à contenir les annotations Swagger
}
