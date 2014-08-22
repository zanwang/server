angular = require 'angular'

angular.module('app').factory 'Auth', ($localStorage) ->
  auth = {}
  auth.token = $localStorage.token

  auth.create = (token) ->
    auth.token = token
    $localStorage.token = token

  auth.destroy = ->
    auth.token = null
    delete $localStorage.token

  auth