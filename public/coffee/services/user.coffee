angular = require 'angular'

angular.module('app').factory 'User', ($resource, Auth, API_BASE_URL, $cacheFactory, $rootScope) ->
  httpCache = $cacheFactory.get '$http'

  $resource API_BASE_URL + 'users/:user_id',
    user_id: if Auth.token then Auth.token.user_id else ''
  ,
    get:
      method: 'GET'
      cache: true
    create:
      method: 'POST'
      url: API_BASE_URL + 'users'
    update:
      method: 'PUT'
      interceptor:
        response: (res) ->
          data = res.resource

          if res.status is 200 and data
            httpCache.put API_BASE_URL + 'users/' + Auth.token.user_id, data
            $rootScope.$emit 'userUpdated', data

          data
