package presentation

import (
	"fmt"

	accountDelivery "disbursement/app/account/delivery"
	accountService "disbursement/app/account/service"
	accountUsecase "disbursement/app/account/usecase"

	disbursementDelivery "disbursement/app/disbursement/delivery"
	disbursementRepository "disbursement/app/disbursement/repository"
	disbursementService "disbursement/app/disbursement/service"
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

	dummyBankServiceClient := infrastructure.NewHTTPClient(config.MOODIEDAM_BANK_API_URL)

	// redsync for distributed mutual exclusion
	// pool := goredis.NewPool(&redisClient)
	// mutexProvider := redsync.New(pool)

	repositories := domain.Repositories{
		DisbursementRepository: disbursementRepository.NewDisbursementRepository(postgresDb),
	}

	services := domain.Services{
		AccountService:      accountService.NewAccountService(dummyBankServiceClient),
		DisbursementService: disbursementService.NewDisbursementService(dummyBankServiceClient),
	}

	usecases := domain.Usecases{
		AccountUsecase:      accountUsecase.NewAccountUsecase(services),
		DisbursementUsecase: disbursementUsecase.NewDisbursementUsecase(repositories, services),
	}

	httpRequestMiddleware := infrastructure.NewHttpRequestMiddleware(config.API_KEY)

	accountDelivery.SetAccountHandler(router, usecases)
	disbursementDelivery.SetDisbursementHandler(router, usecases, httpRequestMiddleware)

	return router
}

func StopServer() {

}
