angular = require 'angular'

angular.module('app').factory 'Domain', ($resource, Auth, API_BASE_URL) ->
  $resource API_BASE_URL + 'domains/:id',
    id: '@id'
  ,
    list:
      method: 'GET'
      url: API_BASE_URL + 'users/:user_id/domains'
      isArray: true
      params:
        user_id: if Auth.token then Auth.token.user_id else ''
    create:
      method: 'POST'
      url: API_BASE_URL + 'users/:user_id/domains'
      params:
        user_id: if Auth.token then Auth.token.user_id else ''