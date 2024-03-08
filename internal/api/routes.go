package api

// Routes applies all routing/endpoint configurations.
func (a *Api) Routes() {
	groupV1 := a.app.Group("/v1")

	docs := groupV1.Group("/docs")
	docs.Name("Get docs").Get("", GetDocs)

	groupCategory := groupV1.Group("/category")

	groupCategory.Name("List Categories").Get("", a.ListCategories)
	groupCategory.Name("Get Category").Get("", a.GetCategory)
	groupCategory.Name("Create Category").Get("", a.CreateCategory)
	groupCategory.Name("Update Categories").Get("", a.UpdateCategory)
	groupCategory.Name("Delete Category").Get("", a.DeleteCategory)
}
