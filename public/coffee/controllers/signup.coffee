angular = require 'angular'

angular.module('app').controller 'SignupCtrl', ($scope, User, $state) ->
  $scope.submitted = false
  $scope.submitting = false
  $scope.success = false
  $scope.user = new User()

  $scope.submit = ->
    return if $scope.submitting

    $scope.submitted = true
    return if $scope.signupForm.$invalid

    $scope.submitting = true

    $scope.user.$create().then (data) ->
      $scope.success = true
    , (err) ->
      $scope.errors = err.data.errors
      $scope.submitting = false