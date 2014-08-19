angular = require 'angular'

angular.module('app').controller 'PasswordCtrl', ($scope, $state) ->
  $scope.submitted = false
  $scope.success = false

  $scope.submit = ->
    $scope.submitted = true
    return if $scope.passwordForm.$invalid