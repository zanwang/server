angular = require 'angular'

angular.module('app').factory 'Auth', ($cookieStore) ->
  auth = {}
  auth.token = $cookieStore.get 'token'

  auth.create = (token) ->
    auth.token = token
    $cookieStore.put 'token', token

  auth.destroy = ->
    auth.token = null
    $cookieStore.remove 'token'

  auth