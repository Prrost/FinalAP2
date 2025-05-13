package useCase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
	"time"
	"user-service/Storage"
	"user-service/Storage/Sqlite"
	"user-service/Storage/cache"
	"user-service/config"
	"user-service/domain"
)

type UseCase struct {
	storage Storage.Storage
	redis   *redis.Client
	cfg     *config.Config
}

func NewUseCase(storage Storage.Storage, cfg *config.Config) *UseCase {
	rdb := cache.NewRedis()

	return &UseCase{
		storage: storage,
		redis:   rdb,
		cfg:     cfg,
	}
}

func (uc *UseCase) CreateUserAdmin(user domain.User) (domain.User, error) {
	const op = "CreateUserAdmin"

	//validation
	val := validator.New()
	err := val.Struct(user)
	if err != nil {
		log.Printf("[%s] validation failed: %v", op, err)
		return domain.User{}, status.Error(codes.InvalidArgument, err.Error())
	}

	//hashing
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[%s] hashing failed: %v", op, err)
		return domain.User{}, status.Error(codes.Internal, err.Error())
	}

	//info that will be stored
	userToStore := domain.User{
		Email:    user.Email,
		Password: string(hash),
		IsAdmin:  user.IsAdmin,
	}

	//storing
	createdUser, err := uc.storage.CreateUserAdmin(userToStore)
	if err != nil {
		if errors.Is(err, Sqlite.ErrUserAlreadyExists) {
			log.Printf("[%s] user is already exist, not created: %v", op, err)
			return domain.User{}, status.Error(codes.AlreadyExists, err.Error())
		}
		log.Printf("[%s] internal DB error: %v", op, err)
		return domain.User{}, status.Error(codes.Internal, err.Error())
	}
	log.Printf("[%s] user created: %s, id: %d", op, user.Email, createdUser.ID)
	return createdUser, nil
}

func (uc *UseCase) CreateUser(user domain.User) (domain.User, error) {
	const op = "CreateUser"

	//validation
	val := validator.New()
	err := val.Struct(user)
	if err != nil {
		log.Printf("[%s] validation failed: %v", op, err)
		return domain.User{}, status.Error(codes.InvalidArgument, err.Error())
	}

	//info that will be stored
	userToStore := domain.User{
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	}

	//storing
	createdUser, err := uc.storage.CreateUser(userToStore)
	if err != nil {
		if errors.Is(err, Sqlite.ErrUserAlreadyExists) {
			log.Printf("[%s] user is already exist, not created: %v", op, err)
			return domain.User{}, status.Error(codes.AlreadyExists, err.Error())
		}
		log.Printf("[%s] internal DB error: %v", op, err)
		return domain.User{}, status.Error(codes.Internal, err.Error())
	}
	log.Printf("[%s] user created: %s, id: %d", op, user.Email, createdUser.ID)
	return createdUser, nil
}

func (uc *UseCase) LoginUser(user domain.User) (string, error) {
	const op = "LoginUser"

	//validating
	val := validator.New()
	err := val.Struct(user)
	if err != nil {
		log.Printf("[%s] validation failed: %v", op, err)
		return "error", status.Error(codes.InvalidArgument, err.Error())
	}

	//taking user from storage
	userStorage, err := uc.storage.GetUserByEmail(user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[%s] user not found in db: %v", op, err)
			return "error", status.Error(codes.NotFound, err.Error())
		}
		log.Printf("[%s] internal DB error: %v", op, err)
		return "error", status.Error(codes.Internal, err.Error())
	}

	if userStorage.IsAdmin == false {
		return "error", status.Error(codes.PermissionDenied, "user is not admin")
	}

	//compare passwords
	log.Printf(userStorage.Password, user.Password)
	err = bcrypt.CompareHashAndPassword([]byte(userStorage.Password), []byte(user.Password))
	if err != nil {
		log.Printf("[%s] vrong credentials: %v", op, err)
		return "error", status.Error(codes.InvalidArgument, err.Error())
	}

	//JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"ID":    userStorage.ID,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	//signing the token
	secretKey := []byte(uc.cfg.JWTSecret)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("[%s] cant sign JWT: %v", op, err)
		return "error", status.Error(codes.Internal, err.Error())
	}

	log.Printf("[%s] registration successful", op)
	return tokenString, nil
}

func (uc *UseCase) GetUserByID(id string) (domain.User, error) {
	const op = "getUserByID"
	var user domain.User

	ctx := context.Background()

	//Ascii to Int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return domain.User{}, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := uc.redis.Get(ctx, id).Result()
	if err == nil {
		err = json.Unmarshal([]byte(res), &user)
		if err != nil {
			log.Printf("[%s] error with cache json unmarshal: %v", op, err)
			return domain.User{}, status.Error(codes.Internal, err.Error())
		}
		log.Printf("[%s] user profile is found in cache", op)
		return user, nil
	} else {
		if errors.Is(err, redis.Nil) {
			log.Printf("[%s] user not found in cache", op)
		} else {
			log.Printf("[%s] Error with cache", op)
			return domain.User{}, status.Error(codes.Internal, err.Error())
		}
	}

	//find in storage
	user, err = uc.storage.GetUserByID(idInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[%s] user not found in db: %v", op, err)
			return domain.User{}, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("[%s] internal DB error: %v", op, err)
		return domain.User{}, status.Error(codes.Internal, err.Error())
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		log.Printf("[%s] error with json marshal: %v", op, err)
		return domain.User{}, status.Error(codes.Internal, err.Error())
	}
	err = uc.redis.Set(ctx, id, userJson, time.Hour).Err()
	if err != nil {
		log.Printf("[%s] error with cache set: %v", op, err)
		return domain.User{}, status.Error(codes.Internal, err.Error())
	}
	log.Printf("[%s] user profile is created in cache", op)

	log.Printf("[%s] user profile is found successfully", op)
	return user, nil
}

func (uc *UseCase) GetUserByEmail(email string) (domain.User, error) {
	const op = "getUserByID"
	ctx := context.Background()
	user := domain.User{}

	res, err := uc.redis.Get(ctx, email).Result()
	if err == nil {
		err = json.Unmarshal([]byte(res), &user)
		if err != nil {
			log.Printf("[%s] error with cache json unmarshal: %v", op, err)
			return domain.User{}, status.Error(codes.Internal, err.Error())
		}
		log.Printf("[%s] user profile is found in cache", op)
		return user, nil
	} else {
		if errors.Is(err, redis.Nil) {
			log.Printf("[%s] user not found in cache", op)
		} else {
			log.Printf("[%s] Error with cache", op)
			return domain.User{}, status.Error(codes.Internal, err.Error())
		}
	}

	//find in storage
	user, err = uc.storage.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[%s] user not found in db: %v", op, err)
			return domain.User{}, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("[%s] internal DB error: %v", op, err)
		return domain.User{}, status.Error(codes.Internal, err.Error())
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		log.Printf("[%s] error with json marshal: %v", op, err)
		return domain.User{}, status.Error(codes.Internal, err.Error())
	}
	err = uc.redis.Set(ctx, email, userJson, time.Hour).Err()
	if err != nil {
		log.Printf("[%s] error with cache set: %v", op, err)
		return domain.User{}, status.Error(codes.Internal, err.Error())
	}
	log.Printf("[%s] user profile is created in cache", op)

	log.Printf("[%s] user profile is found successfully", op)
	return user, nil
}
