package system

import (
	"github.com/gofiber/fiber"
	healthHttp "github.com/thomasvvugt/fiber-boilerplate/routes/system/http"
	_ "github.com/thomasvvugt/fiber-boilerplate/routes/system/mongo"
	healthMySql "github.com/thomasvvugt/fiber-boilerplate/routes/system/mysql"
	healthRabbit "github.com/thomasvvugt/fiber-boilerplate/routes/system/rabbitmq"
	healthRedis "github.com/thomasvvugt/fiber-boilerplate/routes/system/redis"
	"time"
)

func CheckDB(c *fiber.Ctx) {

	// http health check example
	RegisterHealthCheck(Config{
		Name:      "http-check",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: healthHttp.New(healthHttp.Config{
			URL: `http://example.com`,
		}),
	})

	// mysql health check example
	RegisterHealthCheck(Config{
		Name:      "mysql-check",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: healthMySql.New(healthMySql.Config{
			DSN: `root:root@tcp(0.0.0.0:3306)/casbin?charset=utf8`,
		}),
	})

	RegisterHealthCheck(Config{
		Name:      "rabbit-aliveness-check",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: healthRabbit.New(healthRabbit.Config{
			DSN: `amqp://guest:guest@localhost:5672/`,
		}),
	})

	/*RegisterHealthCheck(Config{
		Name:      "mongo-aliveness-check",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: healthMongo.New(healthMongo.Config{
			DSN: `mongodb://localhost:27017`,
		}),
	})*/
	RegisterHealthCheck(Config{
		Name:      "redis-aliveness-check",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: healthRedis.New(healthRedis.Config{
			DSN: `redis://localhost:6379`,
		}),
	})
	HealthInfo(c)
}
