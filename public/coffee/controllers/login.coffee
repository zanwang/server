angular = require 'angular'

angular.module('app').controller 'LoginCtrl', ($scope, Token, $state, Auth) ->
  $scope.submitted = false
  $scope.submitting = false
  $scope.token = new Token()

  $scope.submit = ->
    return if $scope.submitting

    $scope.submitted = true
    return if $scope.loginForm.$invalid

    $scope.submitting = true

    $scope.token.$save().then (token) ->
      Auth.create token
      $state.go 'app.domains'
    , (err) ->
      $scope.errors = err.data.errors
      $scope.submitting = false