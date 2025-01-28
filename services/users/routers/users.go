package routers

import (
	"net/http"

	"github.com/parrothacker1/Solvelt/users/handlers"
	"github.com/parrothacker1/Solvelt/users/middlewares"
)

var UserRouter *Router 

func init() {
  UserRouter = NewRouter()
  UserRouter.Handle(http.MethodPost,"/",handlers.CreateUser)
  UserRouter.Handle(http.MethodPut,"/",handlers.UpdateUser,middlewares.AuthMiddleware)
  UserRouter.Handle(http.MethodGet,"/me",handlers.GetUser,middlewares.AuthMiddleware)
  UserRouter.Handle(http.MethodDelete,"/",handlers.DeleteUser,middlewares.AuthMiddleware)
  UserRouter.Handle(http.MethodPost,"/login",handlers.LoginUser)
  UserRouter.Handle(http.MethodPost,"/reset",handlers.ResetPassword,middlewares.AuthMiddleware)
}
