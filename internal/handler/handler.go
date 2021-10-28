package handler

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"taskESchop/internal/models"
	"taskESchop/internal/service"
)

type Handler struct {
	userService service.UserService
	managerService service.ManagerService
	authorizationService service.AuthorizationService
}

func NewHandler(userService service.UserService, managerService service.ManagerService, authorizationService service.AuthorizationService) *Handler {
	return &Handler{managerService: managerService, userService: userService, authorizationService: authorizationService}
}

//accept username, password, role_ids[]: manager = 1, user = 2
func (h *Handler) SignUp(c echo.Context) error  {
	var req models.User
	if bErr := c.Bind(&req); bErr != nil{
		log.Warnf("bad request: %v", bErr)
		return echo.NewHTTPError(http.StatusBadRequest, bErr)
	}

	res, err := h.authorizationService.SignUp(c.Request().Context(), &req)
	if err != nil{
		return err
	}
	log.Infof("Success signUp response: %v", res)
	return c.JSON(http.StatusCreated, res)
}

func (h *Handler) SignIn(c echo.Context) error  {
	var req models.User
	if bErr := c.Bind(&req); bErr != nil{
		log.Warnf("bad request: %v", bErr)
		return echo.NewHTTPError(http.StatusBadRequest, bErr)
	}
	res, err := h.authorizationService.SignIn(c.Request().Context(), &req)
	if err != nil{
		return err
	}
	token, tErr := CreateToken(res.Id, res.Username, res.RolesName)
	if tErr != nil{
		log.Errorf("Error failed to create token: %v", tErr)
		return echo.NewHTTPError(http.StatusInternalServerError, tErr)
	}
	c.Response().Header().Set("Authorization", token)
	log.Infof("Success signIp response: %v", res)
	return c.JSON(http.StatusOK, res)
}
//accept name, description count
func (h *Handler) CreateNewProduct(c echo.Context) error  {
	if isManager := validateManager(c); !isManager{
		return echo.NewHTTPError(http.StatusForbidden, "not allowed")
	}
	var req models.Product
	if bErr := c.Bind(&req); bErr != nil{
		log.Errorf("bad request")
		return echo.NewHTTPError(http.StatusBadRequest, bErr)
	}
	id, err := h.managerService.CreateNewProduct(c.Request().Context(), &req)
	if err != nil{
		return err
	}
	log.Infof("Product %v success created", id)
	return c.JSON(http.StatusOK, id)
}
//accept all fields of model
func (h *Handler) UpdateProduct(c echo.Context) error {
	if isManager := validateManager(c); !isManager{
		return echo.NewHTTPError(http.StatusForbidden, "not allowed")
	}
	var req models.Product
	if bErr := c.Bind(&req); bErr != nil{
		log.Warnf("bad request: %v", bErr)
		return echo.NewHTTPError(http.StatusBadRequest, bErr)
	}
	err := h.managerService.UpdateProduct(c.Request().Context(), &req)
	if err != nil{
		return err
	}
	log.Infof("Product %v success updated", req.Id)
	return nil
}

//empty request with token
func (h *Handler) GetAllCarts(c echo.Context) error {
	if isManager := validateManager(c); !isManager{
		return echo.NewHTTPError(http.StatusForbidden, "not allowed")
	}
	res, err := h.managerService.GetAllCarts(c.Request().Context())
	if err != nil{
		return err
	}
	log.Infof("Success response: %v", res)
	return c.JSON(http.StatusOK, res)
}


func (h *Handler) AddProductToCart(c echo.Context) error {
	if isUser := validateUser(c); !isUser{
		return echo.NewHTTPError(http.StatusForbidden, "not allowed")
	}
	id := c.Request().Context().Value("user_id").(float64)
	var req models.Product
	if bErr := c.Bind(&req); bErr != nil{
		log.Warnf("bad request: %v", bErr)
		return echo.NewHTTPError(http.StatusBadRequest, bErr)
	}
	err := h.userService.AddProductToCart(c.Request().Context(), &req, int(id))
	if err != nil{
		return err
	}
	return nil
}

func (h *Handler) DeleteFromCart(c echo.Context) error {
	if isUser := validateUser(c); !isUser{
		return echo.NewHTTPError(http.StatusForbidden, "not allowed")
	}
	idStr := c.Param("id")
	id, idErr := strconv.Atoi(idStr)
	if idErr != nil{
		log.Warnf("bad request: %v", idErr)
		return echo.NewHTTPError(http.StatusBadRequest, idErr)
	}
	err := h.userService.DeleteFromCart(c.Request().Context(), id)
	if err != nil{
		return err
	}
	return nil
}

func (h *Handler) GetCart(c echo.Context) error {
	if isUser := validateUser(c); !isUser{
		return echo.NewHTTPError(http.StatusForbidden, "not allowed")
	}
	userId := c.Request().Context().Value("user_id").(float64)
	res, err := h.userService.GetCart(c.Request().Context(), int(userId))
	if err != nil{
		return err
	}
	log.Infof("Success response: %v", res)
	return c.JSON(http.StatusOK, res)
}

func CreateToken(userId int, username string, roles []string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["username"] = username
	claims["role"] = roles

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Info("token created for ", username)
	return jwtToken.SignedString([]byte("secret_key"))
}

func validateManager(c echo.Context) bool {
	roles := c.Request().Context().Value("roles").([]string)
	var isManager bool
	for _, v := range roles{
		if v == "manager"{
			isManager = true
			break
		}
	}
	return isManager
}

func validateUser(c echo.Context) bool {
	roles := c.Request().Context().Value("roles").([]string)
	var isManager bool
	for _, v := range roles{
		if v == "user"{
			isManager = true
			break
		}
	}
	return isManager
}
