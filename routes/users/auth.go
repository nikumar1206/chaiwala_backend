package users

import (
	"errors"
	"log/slog"
	"net/http"

	"ChaiwalaBackend/clients/jwt"
	"ChaiwalaBackend/db"
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
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		usr, err := dbConn.CreateUser(c.Context(), db.CreateUserParams{
			PasswordHash: string(hash),
			Email:        u.Email,
		})
		if err != nil {

			var e *pgconn.PgError
			if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
				c.Status(409)
				return c.JSON(common.Error{
					Message:   "Username already taken.",
					RequestId: c.GetRespHeader("X-Request-ID"),
				})
			}
			c.Status(500)
			return c.JSON(common.Error{
				Message:   "User could not be created",
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		at, rt, exp, err := jwt.GenerateTokens(usr.Email, usr.ID)
		if err != nil {
			c.Status(500)
			return c.JSON(common.Error{
				Message:   "Could not generate a JWT",
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
		slog.InfoContext(c.Context(), "Received a request to loginUser")
		u := new(LoginUser)
		if err := c.Bind().JSON(u); err != nil {
			return err
		}
		usr, err := dbConn.GetUserByEmail(c.Context(), u.Email)
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			c.Status(http.StatusNotFound)
			return common.SendErrorResponse(c, http.StatusNotFound, "User not found")
		}

		err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(u.Password))
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return common.SendErrorResponse(c, http.StatusUnauthorized, "Incorrect password")
		}

		at, rt, exp, err := jwt.GenerateTokens(usr.Email, usr.ID)
		if err != nil {
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Could not generate a JWT")
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
		claims := c.Locals("claims").(jwt.Claims)

		// todo(nick): need to revoke previous Refresh
		newAccess, newRefresh, exp, err := jwt.GenerateTokens(claims.Email, claims.UserID)
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Could not generate JWT")
		}

		// maybe return user info?
		slog.InfoContext(c.Context(), "success")
		return c.Status(200).JSON(GeneratedJWTResponse{
			AccessToken:  newAccess,
			RefreshToken: newRefresh,
			ExpiresIn:    exp.UnixMilli(),
			TokenType:    "Bearer",
		})
	}
}
