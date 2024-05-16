package api

// ListUsers fetches the users
//
// @Id ListUsers
// @Summary Get Users
// @Description Returns the users
// @Tags UsersService
// @Accept application/json
// @Produce application/json
// @Param filters query model.UserFilters false "User filters"
// @Success 200 {object} model.PaginatedUser
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/user [get]
//func (a *Api) ListUsers(ctx *fiber.Ctx) error {
//	filter := model.UserFilters{
//		UserIsActive: []int{1},
//	}
//
//	if err := ctx.QueryParser(&filter); err != nil {
//		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
//	}
//
//	if err := filter.Validate(); err != nil {
//		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
//	}
//	filter.SetPaginationDefaults()
//
//	users, err := a.cfg.UserService.ListUsers(ctx.Context(), &filter)
//	return a.WriteResponse(ctx, http.StatusOK, users, err)
//}
