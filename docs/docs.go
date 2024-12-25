// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "login user by credentials",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "user to login",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.AuthResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "logout user with token's validation",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout user",
                "responses": {
                    "200": {
                        "description": "Token is valid",
                        "schema": {
                            "type": "nil"
                        }
                    },
                    "400": {
                        "description": "Invalid token is sent",
                        "schema": {
                            "type": "nil"
                        }
                    },
                    "401": {
                        "description": "User is unauthorized",
                        "schema": {
                            "type": "nil"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "register user by credentials",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "user to register",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entity.AuthResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/password/request-reset": {
            "post": {
                "description": "recover your password by email code",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recovery"
                ],
                "summary": "Recover password",
                "parameters": [
                    {
                        "description": "email of the user",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.inputRecovery"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code was sent",
                        "schema": {
                            "type": "nil"
                        }
                    },
                    "400": {
                        "description": "invalid email",
                        "schema": {
                            "type": "nil"
                        }
                    }
                }
            }
        },
        "/password/reset": {
            "post": {
                "description": "update password for registered user",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "recovery"
                ],
                "summary": "Reset Password",
                "parameters": [
                    {
                        "description": "password and email of the user for recovery",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.passwordInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "password is valid",
                        "schema": {
                            "type": "nil"
                        }
                    },
                    "400": {
                        "description": "invalid code",
                        "schema": {
                            "type": "nil"
                        }
                    }
                }
            }
        },
        "/password/validate-code": {
            "post": {
                "description": "compare passed code with the saved one",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recovery"
                ],
                "summary": "Validate Code",
                "parameters": [
                    {
                        "description": "code and email of the user for recovery",
                        "name": "code",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.codeInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code is valid",
                        "schema": {
                            "type": "nil"
                        }
                    },
                    "400": {
                        "description": "invalid code",
                        "schema": {
                            "type": "nil"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.AuthResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "entity.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "entity.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user_1234@gmail.com"
                },
                "id": {
                    "type": "string",
                    "example": "1234"
                },
                "password": {
                    "type": "string",
                    "example": "pass_1234"
                },
                "phone": {
                    "type": "string",
                    "example": "89178298123"
                },
                "username": {
                    "type": "string",
                    "example": "user_1234"
                }
            }
        },
        "handlers.codeInput": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                }
            }
        },
        "handlers.inputRecovery": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "handlers.passwordInput": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Auth Service",
	Description:      "This is the auth services of course project",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
