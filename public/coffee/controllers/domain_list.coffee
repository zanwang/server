angular = require 'angular'

angular.module('app').controller 'DomainListCtrl', ($scope, Domain) ->
  $scope.loading = true
  $scope.domains = []
  $scope.submitted = false
  $scope.submitting = false
  $scope.domain = new Domain()

  $scope.loadingStyle =
    width: 2
    lines: 8
    length: 5
    radius: 5

  $scope.create = ->
    return if $scope.submitting

    $scope.submitting = true
    $scope.submitted = true
    return if $scope.domainForm.$invalid

    $scope.domain.$create().then (data) ->
      $scope.domains.push data
      $scope.domain = new Domain()
      $scope.submitted = false
      $scope.submitting = false
      $scope.errors = {}
      $scope.domainForm.$setPristine()
    , (err) ->
      $scope.submitting = false
      $scope.errors = err.data.errors

  # Load domain list
  Domain.list().$promise.then (data) ->
    $scope.domains = data
    $scope.loading = false