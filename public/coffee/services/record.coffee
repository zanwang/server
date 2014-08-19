angular = require 'angular'

angular.module('app').factory 'Record', ($resource, API_BASE_URL) ->
  $resource API_BASE_URL + 'records/:id',
    id: '@id'
  ,
    list:
      method: 'GET'
      url: API_BASE_URL + 'domains/:domain_id/records'
      isArray: true
    create:
      method: 'POST'
      url: API_BASE_URL + 'domains/:domain_id/records'
      params:
        domain_id: '@domain_id'
    update:
      method: 'PUT'