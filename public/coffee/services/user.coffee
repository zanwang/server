angular = require 'angular'

angular.module('app').factory 'User', (Restangular) ->
  Restangular.service 'users'
