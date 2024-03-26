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
        "/artists": {
            "get": {
                "description": "Get list of artists that satisfied conditions in filter",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "artists"
                ],
                "summary": "Get list of artists",
                "parameters": [
                    {
                        "type": "string",
                        "example": "\"Blue Town\"",
                        "description": "name of the artist",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "\"-name\", \"name\"",
                        "description": "criteria for sorting artist-searching results",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 2,
                        "description": "searching page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 10,
                        "description": "searching limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Artist"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "put": {
                "description": "Update information of an artist",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "artists"
                ],
                "summary": "Update information of an artist",
                "parameters": [
                    {
                        "description": "artist information",
                        "name": "artist",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Artist"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Artist"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "description": "Create a new artist",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "artists"
                ],
                "summary": "Create an artist",
                "parameters": [
                    {
                        "description": "Artist Information",
                        "name": "artist",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Artist"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Artist"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/artists/{id}": {
            "get": {
                "description": "Get information of an artist using ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "artists"
                ],
                "summary": "Get information of an artist",
                "parameters": [
                    {
                        "type": "string",
                        "example": "\"3983a1d6-759b-4e5e-b307-7b7e06a05a85\"",
                        "description": "Artist ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Artist"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "delete": {
                "description": "Delete an artist using ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "artists"
                ],
                "summary": "Delete an artist",
                "parameters": [
                    {
                        "type": "string",
                        "example": "\"3983a1d6-759b-4e5e-b307-7b7e06a05a85\"",
                        "description": "Artist ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.DeleteArtistResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/artists/{id}/tracks": {
            "get": {
                "description": "Get top tracks of an artist using ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "artists"
                ],
                "summary": "Get top tracks of an artist",
                "parameters": [
                    {
                        "type": "string",
                        "example": "\"3983a1d6-759b-4e5e-b307-7b7e06a05a85\"",
                        "description": "Artist ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Tracks"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/tracks": {
            "get": {
                "description": "Get information of many tracks satisfied conditions in filter",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tracks"
                ],
                "summary": "Get information of many tracks (advanced)",
                "parameters": [
                    {
                        "type": "string",
                        "example": "\"Blue Town\"",
                        "description": "name of the song",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "\"-namme\", \"name\"",
                        "description": "criteria for sorting track-searching results",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 2,
                        "description": "searching page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 10,
                        "description": "searching limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Tracks"
                        }
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "put": {
                "description": "Update information of a track",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tracks"
                ],
                "summary": "Update information of a track",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Track"
                        }
                    },
                    "400": {
                        "description": "Bad request"
                    }
                }
            },
            "post": {
                "description": "Create a new track",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tracks"
                ],
                "summary": "Create a track",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Track"
                        }
                    },
                    "400": {
                        "description": "Bad request"
                    }
                }
            }
        },
        "/tracks/{id}": {
            "get": {
                "description": "Get information of a track by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tracks"
                ],
                "summary": "Get information of a track",
                "parameters": [
                    {
                        "type": "string",
                        "example": "\"3983a1d6-759b-4e5e-b307-7b7e06a05a85\"",
                        "description": "Track ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Track"
                        }
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            },
            "delete": {
                "description": "Delete a track using ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tracks"
                ],
                "summary": "Delete a track",
                "parameters": [
                    {
                        "type": "string",
                        "example": "\"3983a1d6-759b-4e5e-b307-7b7e06a05a85\"",
                        "description": "Track ID",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.DeleteTrackResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request"
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Artist": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string",
                    "example": "Taylor Swift (born December 13, 1989, West Reading, Pennsylvania, U.S.) is a multitalented singer-songwriter and global superstar who has captivated audiences with her heartfelt lyrics and catchy melodies, solidifying herself as one of the most influential artists in contemporary music."
                },
                "id": {
                    "type": "string",
                    "example": "3983a1d6-759b-4e5e-b307-7b7e06a05a85"
                },
                "name": {
                    "type": "string",
                    "example": "Taylor Swift"
                }
            }
        },
        "model.Track": {
            "type": "object",
            "properties": {
                "artistID": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "3983a1d6-759b-4e5e-b307-7b7e06a05a85"
                    ]
                },
                "id": {
                    "type": "string",
                    "example": "3983a1d6-759b-4e5e-b307-7b7e06a05a85"
                },
                "length": {
                    "type": "integer",
                    "example": 88
                },
                "name": {
                    "type": "string",
                    "example": "Blue Town"
                }
            }
        },
        "model.Tracks": {
            "type": "object",
            "properties": {
                "tracks": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
                    "example": {
                        "key": "value"
                    }
                }
            }
        },
        "response.DeleteArtistResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "type": "string"
                }
            }
        },
        "response.DeleteTrackResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:4040",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Swagger Flotify API",
	Description:      "Spotify API clone",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
