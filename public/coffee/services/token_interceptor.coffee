angular = require 'angular'

angular.module('app').factory 'TokenInterceptor', (Auth, $window, $q) ->
  interceptor = {}

  # Set authorization token
  interceptor.request = (config) ->
    if Auth.token
      config.headers.Authorization = "token #{Auth.token.key}"

    config

  # Handle token error
  interceptor.responseError = (res) ->
    # Unauthorized
    if Auth.token and res.status is 401
      Auth.destroy()
      $window.location.href = '/login'

    $q.reject res

  interceptor