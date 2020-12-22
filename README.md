This is the go boilerplate on the top of fiber web framework.

With simple setup, you can use many features out of the box:

The features include:

* Basic Auth with Login,Register
* Email confirmation on Registration
* Role based authorization using Casbin
* File uploads
* UI on Tailwind. Setup ready for VueJS integration
* Laravel mix for UI
* Payment processing via PayPal
* Logging via Zerolog with file rotation
* MySQL with GORM
* HTTP Client with Retry (with Backoff strategy), Throttle and Timeout
* RabbitMQ Integration for Queue Processing
* REST API Authentication with JWT
* APP and API Separation based on JWT Token
* REST based basic auth
* Use of Redis for Cache and Session
* Hot Reload with Air
* Flash Message based on cookies
* Easy Config Settings based on .env
* Setup for Docker

# Installation
* Clone the repo `git clone https://github.com/sujit-baniya/fiber-boilerplate.git`
* Make sure you have installed: Redis, MySQL or Postgres, RabbitMQ (Optional)
* Copy .env.sample to .env
* If you're not using RabbitMQ, comment line:50 and line:52 on main.go
* To build the frontend, install nodejs
* Then run `npm install` and `npm run prod`
* Your server should be up now


Thanks to following libraries:

* [Fiber](https://github.com/gofiber/fiber/v2)
* [Xopen](https://github.com/brentp/xopen)
* [Gorm](https://github.com/go-gorm/gorm)
* [Zerolog](https://github.com/edersohe/zflogger)
* [Jwt](https://github.com/dgrijalva/jwt-go)
* [Fiber-Casbin](https://github.com/arsmn/fiber-casbin)
* [Air](https://github.com/cosmtrek/air)

[![Donate with PayPal](https://raw.githubusercontent.com/itsursujit/fiber-boilerplate/master/paypal-donate-button.png)](https://www.paypal.me/spbaniya)
