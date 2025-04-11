package users

import (
	"errors"
	"fmt"
	"net/http"

	"ChaiwalaBackend/db"
	"ChaiwalaBackend/jwt"
	common "ChaiwalaBackend/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func BuildAuthRouter(app *fiber.App, dbConn *db.Queries) *fiber.Router {
	userRouter := app.Group("/auth")

	userRouter.Get("", getUser(dbConn))
	userRouter.Post("/register", registerUser(dbConn))
	userRouter.Post("/login", loginUser(dbConn))
	userRouter.Post("/refresh", refreshRoute(dbConn))

	return &userRouter
}

func registerUser(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		u := new(RegisterUser)
		if err := c.Bind().JSON(u); err != nil {
			return err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0)
		if err != nil {
			c.Status(500)
			return c.JSON(common.Error{
				Message:   "User could not be created",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		usr, err := dbConn.CreateUser(c.Context(), db.CreateUserParams{
			Username:     u.Username,
			PasswordHash: string(hash),
			Email:        u.Username + "@gmail.com",
		})
		if err != nil {

			var e *pgconn.PgError
			if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
				c.Status(409)
				return c.JSON(common.Error{
					Message:   "Username already taken.",
					Context:   err.Error(),
					RequestId: c.GetRespHeader("X-Request-ID"),
				})
			}
			c.Status(500)
			return c.JSON(common.Error{
				Message:   "User could not be created",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		at, rt, exp, err := jwt.GenerateTokens(usr.Username, usr.ID)
		if err != nil {
			c.Status(500)
			return c.JSON(common.Error{
				Message:   "Could not generate a JWT",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		c.Status(200)
		return c.JSON(GeneratedJWTResponse{
			AccessToken:  at,
			RefreshToken: rt,
			ExpiresIn:    exp.UnixMilli(),
			TokenType:    "Bearer",
		})
	}
}

func loginUser(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		u := new(LoginUser)
		if err := c.Bind().JSON(u); err != nil {
			return err
		}
		usr, err := dbConn.GetUserByUsername(c.Context(), u.Username)
		if err != nil {
			fmt.Print(err)
			c.Status(http.StatusNotFound)
			return c.JSON(common.Error{
				Message:   "User not found.",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(u.Password))
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return c.JSON(common.Error{
				Message:   "Incorrect password",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		at, rt, exp, err := jwt.GenerateTokens(usr.Username, usr.ID)
		if err != nil {
			c.Status(500)
			return c.JSON(common.Error{
				Message:   "Could not generate a JWT",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		return c.JSON(
			LoginUserResponse{
				Token: GeneratedJWTResponse{
					AccessToken:  at,
					RefreshToken: rt,
					ExpiresIn:    exp.UnixMilli(),
					TokenType:    "Bearer",
				},
				User: usr,
			},
		)
	}
}

func refreshRoute(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		req := new(RefreshTokenRequest)
		if err := c.Bind().JSON(req); err != nil {
			return err
		}

		claims, err := jwt.ValidateToken(req.RefreshToken)
		if err != nil {
			c.Status(401)
			return c.JSON(common.Error{
				Message:   "Not a Valid Token",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})

		}

		// todo(nick): need to revoke previous Refresh
		newAccess, newRefresh, exp, err := jwt.GenerateTokens(claims.Username, claims.UserID)
		if err != nil {
			c.Status(500)
			return c.JSON(common.Error{
				Message:   "Could not generate a JWT",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		c.Status(200)
		return c.JSON(GeneratedJWTResponse{
			AccessToken:  newAccess,
			RefreshToken: newRefresh,
			ExpiresIn:    exp.UnixMilli(),
			TokenType:    "Bearer",
		})
	}
}
