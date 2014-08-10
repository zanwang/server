angular = require 'angular'

angular.module('app').factory 'Token', (Restangular) ->
  Restangular.service 'tokens'
