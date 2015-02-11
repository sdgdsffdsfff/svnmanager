define([
'kernel',
'angular',
'./module',
'react',
'ui/Flyout',
'components/ui/Wrap',
'components/ui/Mod',
'service/SvnService'
],
function( core, ng, directive, React, Flyout, Wrap, Mod ){
    directive
        .directive('svnUpdate', function (SvnService) {

            var flyout = new Flyout(Wrap(Mod, {
                title: "Upgrade.."
            }));

            return {
                controller: function ($scope) {
                    $scope.formatTime = function (str) {
                        return moment(str).format("YYYY-MM-DD HH:mm:ss")
                    };

                    $scope.svnUpdate = function () {
                        return SvnService.svnup().then(function (data) {
                            $scope.upFileList = data.List;
                            $scope.deployEnable = true;
                            ng.extend($scope.version, data.Version);
                        }, function (data) {
                            console.log(data)
                        })
                    }
                },
                link: function (scope, elem) {
                    elem.click(function () {
                        flyout.show(this, 'bottom', 'right');
                    });
                }
            }
        })
        .directive('svnDeploy', function (SvnService) {
            return {
                controller: function ($scope) {
                    $scope.deploy = function (ids) {
                        return SvnService.deploy(ids).then(function (data) {
                            console.log(data)
                        })
                    }
                },
                link: function (scope, elem) {
                    elem.click(function () {
                        scope.deploy()
                    })
                }
            }
        })
});