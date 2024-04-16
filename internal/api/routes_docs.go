package api

import "github.com/gofiber/fiber/v2"

func GetDocs(ctx *fiber.Ctx) error {
	return ctx.Render("index", fiber.Map{
		"Title":   "Hello, World!",
		"baseUrl": "www.google.com",
	})
}

const (
	allIReallyNeedIsYou = `
		<!DOCTYPE html>
			<html>
			<head>
				<title>Simple Hello World API Documentation</title>
				<!-- Load ReDoc -->
				<script src="https://cdn.jsdelivr.net/npm/redoc/bundles/redoc.standalone.js"></script>
			</head>
			<body>
				<!-- Container where ReDoc will render the API documentation -->
				<div id="redoc-container"></div>
				
				<!-- Embedded YAML OpenAPI specification for a simple Hello World API -->
				<script id="embedded-yaml" type="application/yaml">
			openapi: 3.0.0
			info:
			  title: Hello World API
			  description: A simple Hello World API
			  version: 1.0.0
			servers:
			  - url: 'http://example.com/api'
			paths:
			  /hello:
				get:
				  summary: Returns a Hello World message
				  responses:
					'200':
					  description: A simple Hello World message
					  content:
						application/json:
						  schema:
							type: object
							properties:
							  message:
								type: string
								example: Hello, World!
				</script>
			
				<script>
				// Function to load and render the embedded YAML with ReDoc
				function renderDocumentation() {
					// Extract the YAML content from the <script> tag
					var yamlContent = document.getElementById('embedded-yaml').textContent;
			
					// Initialize ReDoc
					Redoc.init(yamlContent, {}, document.getElementById('redoc-container'));
				}
			
				// Render the documentation after the page has loaded
				window.onload = renderDocumentation;
				</script>
			</body>
			</html>

	`
)
