angular = require 'angular'

angular.module('app').controller 'RecordCtrl', ($scope) ->
  $scope.submitted = false
  $scope.editing = false
  $scope.submitting = false

  $scope.edit = ->
    $scope.editing = true
    $scope.backup = angular.copy $scope.record

  $scope.cancel = ->
    $scope.editing = false
    $scope.record = $scope.backup

  $scope.update = ->
    return if $scope.submitting

    $scope.submitted = true
    return if $scope.editForm.$invalid

    $scope.submitting = true

    $scope.record.$update().then (data) ->
      $scope.record = data
      $scope.submitted = false
      $scope.submitting = false
      $scope.editing = false
      $scope.errors = {}
      $scope.recordForm.$setPristine()
    , (err) ->
      $scope.errors = err.data.errors
      $scope.submitting = false

  $scope.delete = ->
    if confirm 'Do you want to delete this record?'
      $scope.record.$delete().then ->
        $scope.$parent.$emit 'recordDeleted', $scope.record.id