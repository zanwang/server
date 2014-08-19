angular = require 'angular'

angular.module('app').factory 'User', ($resource, Auth, API_BASE_URL, $cacheFactory, $rootScope) ->
  httpCache = $cacheFactory.get '$http'

  $resource API_BASE_URL + 'users/:user_id',
    user_id: if Auth.token then Auth.token.user_id else ''
  ,
    get:
      method: 'GET'
      cache: true
    update:
      method: 'PUT'
      interceptor:
        response: (res) ->
          if res.status is 200 and res.data
            httpCache.put API_BASE_URL + 'users/' + Auth.token.user_id, res.data
            $rootScope.$emit 'userUpdated', res.data

          res
