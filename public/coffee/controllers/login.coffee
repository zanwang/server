angular = require 'angular'

angular.module('app').controller 'LoginCtrl', ($scope, Token) ->
  $scope.login = ->
    form = $scope.loginForm
