angular = require 'angular'

angular.module('app').controller 'DomainCtrl', ($scope, Record) ->
  $scope.loading = false
  $scope.submitted = false
  $scope.submitting = false
  $scope.expanded = false
  $scope.loaded = false
  $scope.records = []
  $scope.recordTypes = ['A', 'CNAME', 'MX', 'TXT', 'SPF', 'AAAA', 'NS', 'LOC']

  $scope.recordTTL = [
    {id: 1, name: 'Automatic'}
    {id: 300, name: '5 mins'}
    {id: 600, name: '10 mins'}
    {id: 900, name: '15 mins'}
    {id: 1800, name: '30 mins'}
    {id: 3600, name: '1 hour'}
    {id: 7200, name: '2 hours'}
    {id: 18000, name: '5 hours'}
    {id: 43200, name: '12 hours'}
    {id: 86400, name: '1 day'}
  ]

  $scope.recordTypeHint =
    A: 'e.g. 127.0.0.1'
    CNAME: 'e.g. mydomain.com'
    MX: 'e.g. mydomain.com'
    TXT: 'Text record value'
    SPF: 'SPF record value'
    AAAA: 'e.g. ::1'
    NS: 'e.g. a.nameserver.com'
    LOC: 'Loc record value'

  $scope.TTLtxt = {}

  angular.forEach $scope.recordTTL, (ttl) ->
    $scope.TTLtxt[ttl.id] = ttl.name

  newRecord = ->
    record = new Record()
    record.type = 'A'
    record.ttl = 1

    record

  $scope.record = newRecord()

  $scope.show = ->
    $scope.expanded = !$scope.expanded

    if $scope.expanded and !$scope.loaded
      $scope.loading = true

      # Load record list
      Record.list
        domain_id: $scope.domain.id
      .$promise.then (data) ->
        $scope.records = data
        $scope.loading = false
        $scope.loaded = true

    # Clean record data
    unless $scope.expanded
      $scope.record = newRecord()

  $scope.create = ->
    return if $scope.submitting

    $scope.submitting = true
    $scope.submitted = true
    return if $scope.recordForm.$invalid

    $scope.record.domain_id = $scope.domain.id

    $scope.record.$create().then (data) ->
      $scope.records.push data
      $scope.submitted = false
      $scope.submitting = false
      $scope.record = newRecord()
      $scope.errors = {}
      $scope.recordForm.$setPristine()
    , (err) ->
      $scope.errors = err.data.errors
      $scope.submitting = false

  $scope.$on 'recordDeleted', (event, id) ->
    for i, record of $scope.records
      if record.id is id
        $scope.records.splice(i, 1)
        break