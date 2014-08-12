angular = require 'angular'

angular.module('app').controller 'LoginCtrl', ($scope) ->
  $scope.login = ->
    return if $scope.loginForm.$invalid
