// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support"
        },
        "license": {
            "name": "MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/articles": {
            "get": {
                "description": "获取多篇文章",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章"
                ],
                "summary": "Get multiple articles",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "标签ID",
                        "name": "tag_id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "状态",
                        "name": "state",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "每页数量",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "添加文章",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章"
                ],
                "summary": "Add a article",
                "parameters": [
                    {
                        "description": "Article",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/blogArticle.Article"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/articles/{id}": {
            "get": {
                "description": "获取文章",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章"
                ],
                "summary": "Get a single article",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "文章ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "更新文章",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章"
                ],
                "summary": "Update a article",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "文章",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/blogArticle.Article"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "删除文章",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "文章"
                ],
                "summary": "Delete a article",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/tags": {
            "get": {
                "description": "Get multiple article tags",
                "produces": [
                    "application/json"
                ],
                "summary": "GetTags",
                "parameters": [
                    {
                        "type": "string",
                        "description": "标签名称",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "状态",
                        "name": "state",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "每页数量",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Add multiple article tags",
                "produces": [
                    "application/json"
                ],
                "summary": "AddTags",
                "parameters": [
                    {
                        "type": "string",
                        "description": "标签名称",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "状态",
                        "name": "state",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "创建人",
                        "name": "created_by",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/tags/{id}": {
            "put": {
                "description": "Edit multiple article tags",
                "produces": [
                    "application/json"
                ],
                "summary": "EditTags",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "标签名称",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "状态",
                        "name": "state",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "修改人",
                        "name": "modified_by",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "summary": "DeleteTags",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth": {
            "get": {
                "description": "获取用户信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "获取用户信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "密码",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "blogArticle.Article": {
            "type": "object",
            "required": [
                "created_on",
                "id",
                "modified_on"
            ],
            "properties": {
                "content": {
                    "type": "string",
                    "minLength": 1
                },
                "created_by": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                },
                "created_on": {
                    "description": "v0.2.2 前写错了类型",
                    "type": "string"
                },
                "deleted_on": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                },
                "desc": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                },
                "id": {
                    "type": "integer"
                },
                "modified_by": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                },
                "modified_on": {
                    "type": "string"
                },
                "state": {
                    "type": "integer",
                    "enum": [
                        0,
                        1
                    ]
                },
                "tag": {
                    "$ref": "#/definitions/blogTag.Tag"
                },
                "tag_id": {
                    "type": "integer",
                    "minimum": 1
                },
                "title": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                }
            }
        },
        "blogTag.Tag": {
            "type": "object",
            "properties": {
                "created_by": {
                    "type": "string"
                },
                "created_on": {
                    "description": "数据库时间改为varchar了",
                    "type": "string"
                },
                "id": {
                    "type": "integer",
                    "minimum": 1
                },
                "modified_by": {
                    "type": "string"
                },
                "modified_on": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "state": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8088",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "gin-gorm-practice",
	Description:      "gin-gorm-practice",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}