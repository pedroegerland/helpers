package echoswagger

import (
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

func configureSwaggerFunc(swaggerFileRoute string) func(c *Config) {
	return func(c *Config) {
		c.URLs = []string{swaggerFileRoute}
	}
}

func RegisterEchoSwagger(e *echo.Echo, version string) (string, string, error) {
	_, callerPath, _, _ := runtime.Caller(1) //nolint:dogsled
	swaggerFilePath := resolveSwaggerFilePath(callerPath)
	swaggerFileRoute := resolveSwaggerFileRoute()

	swaggerProxyPrefix := os.Getenv("SWAGGER_PROXY_PREFIX")
	if swaggerProxyPrefix != "" {
		file, err := os.ReadFile(swaggerFilePath)
		if err != nil {
			return "", "", err
		}

		var obj any
		err = yaml.Unmarshal(file, &obj)
		if err != nil {
			return "", "", err
		}

		if objMap, ok := obj.(map[string]any); ok {
			objMap["servers"] = []any{map[string]string{"url": swaggerProxyPrefix}}
			if objMapInfo, ok := objMap["info"]; ok {
				if objMapInfoMap, ok := objMapInfo.(map[string]any); ok {
					objMapInfoMap["version"] = version
				}
			} else {
				objMap["info"] = map[string]any{"version": version}
			}

			out, err := yaml.Marshal(objMap)
			if err != nil {
				return "", "", err
			}

			err = os.WriteFile(swaggerFilePath, out, os.ModePerm)
			if err != nil {
				return "", "", err
			}
		}
	}

	e.File("/echo-swagger/echo-swagger.yml", swaggerFilePath)
	e.GET(path.Join("/echo-swagger", "*"), EchoWrapHandler(configureSwaggerFunc(swaggerFileRoute)))

	return swaggerFilePath, swaggerFileRoute, nil
}

func resolveSwaggerFilePath(defaultPath string) string {
	swaggerPath, configured := os.LookupEnv("SWAGGER_PATH")
	if !configured {
		swaggerPath = filepath.Dir(defaultPath)
	}

	sourceFilePath := filepath.Join(swaggerPath, "echo-swagger.yml")
	return sourceFilePath
}

func resolveSwaggerFileRoute() string {
	swaggerProxyPath := os.Getenv("SWAGGER_PROXY_PREFIX")
	fileRoute := filepath.Join(swaggerProxyPath, "/echo-swagger", "echo-swagger.yml")
	return fileRoute
}
