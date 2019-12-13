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

var graphSchema, nodeSchema, edgeShema gojsonschema.JSONLoader
var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	graphSchema = gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s/schema/graph.schema.json", path))
	nodeSchema = gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s/schema/node.schema.json", path))
	edgeShema = gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s/schema/edge.schema.json", path))
}

func validateGraphJSONBody(body json.RawMessage) error {
	loader :=  gojsonschema.NewBytesLoader(body)
	result, err := gojsonschema.Validate(graphSchema, loader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		ve := ValidationError{}
		for _, err := range result.Errors() {
			if err.Field() == gojsonschema.STRING_ROOT_SCHEMA_PROPERTY {
				ve["graph"] = err.Description()
			} else {
				ve[err.Field()] = err.Description()
			}
		}

		return ve
	}

	return nil
}

func graphPostHandler(ectx echo.Context) error {
	var jsonBody json.RawMessage
	if err := ectx.Bind(&jsonBody); err != nil {
		return ectx.JSON(http.StatusBadRequest, HTTPErrorResponse{
			Message: "failed to parse request data as JSON bytes",
			Error:   err,
		})
	}

	if err := validateGraphJSONBody(jsonBody); err != nil {
		return ectx.JSON(http.StatusBadRequest, HTTPErrorResponse{
			Message: "failed to pass validation",
			Error:   err,
		})
	}

	return ectx.JSON(http.StatusOK, HTTPSuccessResponse{
		Message:     "directed acyclic graph is successfully created",
		RequestBody: jsonBody,
	})
}

func main() {
	e := echo.New()
	e.GET("/", func(ectx echo.Context) error {
		return ectx.String(http.StatusOK, "Hello, World!")
	})

	api := e.Group("api/")
	api.POST("dags/", graphPostHandler)
	e.Logger.Fatal(e.Start(":8000"))
}
