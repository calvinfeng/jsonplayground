package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

var animalSchema gojsonschema.JSONLoader
var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	animalSchema = gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s/schema/animal.schema.json", path))
}

func main() {

	e := echo.New()
	e.GET("/", func(ectx echo.Context) error {
		return ectx.String(http.StatusOK, "Hello, World!")
	})

	api := e.Group("api/")
	api.POST("animals/", func(ectx echo.Context) error {
		var jsonBody json.RawMessage
		if err := ectx.Bind(&jsonBody); err != nil {
			return ectx.JSON(http.StatusBadRequest, HTTPErrorResponse{
				Message: "failed to parse request data as JSON bytes",
				Error:   err,
			})
		}

		jsonLoader := gojsonschema.NewBytesLoader(jsonBody)
		result, err := gojsonschema.Validate(animalSchema, jsonLoader)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, HTTPErrorResponse{
				Message: "failed to perform JSON validation",
				Error:   err,
			})
		}

		if !result.Valid() {
			ve := ValidationError{}
			for _, err := range result.Errors() {
				if err.Field() == gojsonschema.STRING_ROOT_SCHEMA_PROPERTY {
					ve["animal"] = err.Description()
				} else {
					ve[err.Field()] = err.Description()
				}
			}
			return ectx.JSON(http.StatusBadRequest, HTTPErrorResponse{
				Message: "data do not pass validation",
				Error:   ve,
			})
		}

		return ectx.JSON(http.StatusOK, HTTPSuccessResponse{
			Message:     "animal is successfully created",
			RequestBody: jsonBody,
		})
	})

	e.Logger.Fatal(e.Start(":8000"))
}
