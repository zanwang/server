angular = require 'angular'

angular.module('app').controller 'SettingsCtrl', ($scope, User, Auth, $state) ->
  $scope.loading = true
  $scope.submitted = false
  $scope.submitting = false

  User.get().$promise.then (user) ->
    $scope.loading = false
    $scope.user = user

  $scope.submit = ->
    return if $scope.submitting

    $scope.submitted = true
    return if $scope.settingForm.$invalid

    $scope.submitting = true

    $scope.user.$update().then (data) ->
      $scope.user = data
      $scope.submitted = false
      $scope.submitting = false
      $scope.errors = {}
      $scope.settingForm.$setPristine()
    , (err) ->
      $scope.errors = err.data.errors
      $scope.submitting = false

  $scope.delete = ->
    if confirm 'Are you sure you want to delete this account?'
      $scope.user.$delete().then ->
        Auth.destroy()
        $state.go 'login'