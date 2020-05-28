# Fiber Boilerplate
A boilerplate for the Fiber web framework


## Configuration
All configuration for your application can be found in the `./config` directory. Various options can be changed depending on your needs such as Database settings, Fiber settings and Fiber Middleware setting such as Logger, Public and Helmet.

These configurations can be found in different files such as `app.yaml`, `fiber.yaml` and `template.yaml`.

Keep in mind if configurations are not set, they default to Fiber's default settings which can be found [here](https://docs.gofiber.io/).


## Routing
Routing examples can be found within the `/routes` directory. Both web and API routes are split, but you can adjust this to your likings.


## Views
Views can be found and edited under `/resources/views` or any other directory to your liking using the `fiber.yaml` and `template.yaml`configuration files.


## Database
We use GORM as an ORM to provide useful features to your models. Please check out their documentation [here](https://github.com/jinzhu/gorm).


## Controllers
Example controllers can be found within the `/app/controllers` directory. You can extend or edit these to your preferences.


## Models
Models are located within the `/app/models` directory.


## Providers
Providers (custom middleware) can be found at `/app/providers`. These providers are not automatically registered.

## Compiling assets
This boilerplate uses [Laravel Mix](https://github.com/JeffreyWay/laravel-mix) as an elegant wrapper around [Webpack](https://github.com/webpack/webpack) (a bundler for javascript and friends).

In order to compile your assets, you must first add them in the `webpack.mix.js` file. Examples of the Laravel Mix API can be found [here](https://laravel-mix.com/docs/5.0/mixjs).

Then you must run either `npm install` or `yarn install` to install the packages required to compile your assets.

Next, run one of the following commands to compile your assets with either `npm` or `yarn`:

```bash
# Run all Mix tasks
npm run dev

# Run all Mix tasks and minify output
yarn run production
# Run all Mix tasks and watch for changes (useful when developing)
yarn run watch
# Run all Mix tasks with hot module replacement
yarn run hot
```

## Docker
You can run your own application using the Docker example image.
To build and run the Docker image, you can use the following commands.

```bash
docker build -t fiber-boilerplate .
docker run --name fiber-boilerplate -p 3000:3000 fiber-boilerplate
```

## Live Reloading (Air)
Example configuration files for [Air](https://github.com/cosmtrek/air) have also been included.
This allows you to live reload your Go application when you change a model, view or controller.

To run Air, use the following commands. Also, check out [Air its documentation](https://github.com/cosmtrek/air) about running the `air` command.
```bash
# Windows
air -c .air.windows.conf
# Linux
air -c .air.linux.conf
```
