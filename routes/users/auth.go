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
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Could not hash password.")
		}

		usr, err := dbConn.CreateUser(c.Context(), db.CreateUserParams{
			PasswordHash: string(hash),
			Email:        u.Email,
		})
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			var e *pgconn.PgError
			if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
				return common.SendErrorResponse(c, http.StatusConflict, "Username already taken")
			}
			return common.SendErrorResponse(c, http.StatusInternalServerError, "User could not be created")
		}

		at, rt, exp, err := jwt.GenerateTokens(usr.Email, usr.ID)
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Could not generate JWT")
		}
		slog.InfoContext(c.Context(), "User created successfully")
		return c.Status(200).JSON(GeneratedJWTResponse{
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
			return common.SendErrorResponse(c, http.StatusNotFound, "User not found")
		}

		err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(u.Password))
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusUnauthorized, "Incorrect password")
		}

		at, rt, exp, err := jwt.GenerateTokens(usr.Email, usr.ID)
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
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
