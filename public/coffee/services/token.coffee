angular = require 'angular'

angular.module('app').factory 'Token', ($resource, API_BASE_URL) ->
  $resource API_BASE_URL + 'tokens'