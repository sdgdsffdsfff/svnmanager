define([
'kernel',
'angular',
'./module',
'ui/Flyout',
'components/ui/upfileList',
'service/SvnService'
],
function( core, ng, directive, Flyout, upfileList ){
    directive
        .directive('svnUpdate', function (SvnService) {
            var flyout = new Flyout( upfileList(), {
                classStyle: 'box'
            });

            var list = [{"Action":3,"Path":"/css.js"},{"Action":3,"Path":"/data.js"},{"Action":1,"Path":"/attributes/css.js"},{"Action":1,"Path":"/attributes/data.js"},{"Action":3,"Path":"/core/attr.js"},{"Action":3,"Path":"/core/init.js"},{"Action":3,"Path":"/core/ready.js"},{"Action":1,"Path":"/effects/init.js"},{"Action":1,"Path":"/effects/ready.js"},{"Action":1,"Path":"/effects/attr.js"}];

            return {
                controller: function ($scope) {
                    $scope.formatTime = function (str) {
                        return moment(str).format("YYYY-MM-DD HH:mm:ss")
                    };

                    $scope.svnUpdate = function () {
                        return SvnService.svnup().then(function (data) {
                            //$scope.upFileList = data.List;
                            //$scope.deployEnable = true;
                            //ng.extend($scope.version, data.Version);
                        }, function (data) {
                            console.log(data)
                        })
                    }
                },
                link: function (scope, elem) {
                    elem.click(function () {

                        flyout.getReact().setList(list);

                        flyout.show(this, 'bottom', 'right');
                        //scope.svnUpdate();
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