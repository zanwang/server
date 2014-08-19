angular = require 'angular'

angular.module('app').controller 'SettingsCtrl', ($scope, User) ->
  $scope.submitted = false

  User.get().$promise.then (user) ->
    $scope.user = user

  $scope.submit = ->
    $scope.submitted = true
    return if $scope.passwordForm.$invalid