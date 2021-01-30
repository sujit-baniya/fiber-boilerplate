This is the go boilerplate on the top of fiber web framework.

With simple setup, you can use many features out of the box.

For details, visit [Documentation](https://sujit-baniya.gitbook.io/fiber-boilerplate/)

The features include:

* Basic Auth with Login,Register
* Email confirmation on Registration
* Role based authorization using Casbin
* File uploads
* UI on Tailwind. Setup ready for VueJS integration
* Laravel mix for UI
* Payment processing via PayPal
* Logging via Phuslu/Log with file rotation
* PostGres or MySQL with GORM V2
* REST API Authentication with JWT
* APP and API Separation based on JWT Token
* REST based basic auth
* Use of Redis for Cache and Session
* Hot Reload with Air
* Flash Message based on cookies
* Easy Config Settings based on .env
* Setup for Docker
* Easy and Almost Zero Downtime Production Deployment with Makefile

# Installation
* Clone the repo `git clone https://github.com/sujit-baniya/fiber-boilerplate.git`
* Make sure you have installed: Redis, MySQL or Postgres
* Copy .env.sample to .env
* To build the frontend, install nodejs
* Then run `npm install` and `npm run prod`
* Your server should be up now

Thanks to following libraries:

* [Fiber](https://github.com/gofiber/fiber/v2)
* [Xopen](https://github.com/brentp/xopen)
* [Gorm](https://github.com/go-gorm/gorm)
* [Phuslu Log](https://github.com/phuslu/log)
* [Jwt](github.com/form3tech-oss/jwt-go)
* [Fiber-Casbin](https://github.com/arsmn/fiber-casbin)
* [Air](https://github.com/cosmtrek/air)