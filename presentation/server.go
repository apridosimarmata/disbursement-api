package presentation

import (
	"fmt"

	accountDelivery "disbursement/app/account/delivery"
	accountService "disbursement/app/account/service"
	accountUsecase "disbursement/app/account/usecase"

	disbursementDelivery "disbursement/app/disbursement/delivery"
	disbursementRepository "disbursement/app/disbursement/repository"
	disbursementUsecase "disbursement/app/disbursement/usecase"

	"disbursement/domain"
	"disbursement/infrastructure"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitServer() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	config := infrastructure.GetConfig()

	fmt.Println(fmt.Sprintf("config got: %v", config))

	postgresDb := infrastructure.NewPostgresConn(config)
	// redisClient := infrastructure.NewRedisClient(ctx, config)
	// cache := infrastructure.NewCache(redisClient)

	dummyBankServiceClient := infrastructure.NewHTTPClient("https://91d6ea63-ae7f-4167-936f-699986c9fa36.mock.pstmn.io/api/v1")

	/***

	{
		"status" : "success",
		"data" : {
			"createdAt": 1712226994,
			"name": "name 1",
			"number": "211833558",
			"id": "1"
		}
	}

	{
		"status" : "error",
		"data" : {
			"message" : "account not found"
		}
	}

	***/

	// redsync for distributed mutual exclusion
	// pool := goredis.NewPool(&redisClient)
	// mutexProvider := redsync.New(pool)

	repositories := domain.Repositories{
		DisbursementRepository: disbursementRepository.NewDisbursementRepository(postgresDb),
	}

	services := domain.Services{
		AccountService: accountService.NewAccountService(dummyBankServiceClient),
	}

	usecases := domain.Usecases{
		AccountUsecase:      accountUsecase.NewAccountUsecase(services),
		DisbursementUsecase: disbursementUsecase.NewDisbursementUsecase(repositories, services),
	}

	accountDelivery.SetAccountHandler(router, usecases)
	disbursementDelivery.SetDisbursementHandler(router, usecases)

	return router
}

func StopServer() {

}
