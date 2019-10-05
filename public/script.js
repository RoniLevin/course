var app = angular.module("FlowApp", ['flow']);
app.config(['flowFactoryProvider', function(flowFactoryProvider) {
  flowFactoryProvider.defaults = {
    target: '/upload',
    chunkSize: 1024 * 1024 * 10,
    maxChunkRetries: 1,
    simultaneousUploads: 1,
    testChunks: false,
    permanentErrors:[404, 500, 501]
  }
}]);

app.controller("FlowCtrl", ["$scope", function($scope) {
  $scope.percentDone = function(file) {
    return (file.progress(false) * 100).toFixed(2).toString() + "%";
  };

  $scope.progress = function(file) {
    return {width: $scope.percentDone(file)};
  };

  $scope.isDone = function(file) {
    return file.isComplete()
  }
}]);
