define([
'kernel',
'angular',
'./module',
'ui/Dialog',
'ui/tips',
'react',
'components/ui/upgradeDialog',
'service/SvnService'
],
function( core, ng, directive, Dialog, tips, React, upgradeDialog ){
    directive
        .directive('svnUpdate', function (SvnService) {
            var dialog;

            var list = [{"Action":3,"Path":"/css.js"},{"Action":3,"Path":"/data.js"},{"Action":1,"Path":"/attributes/css.js"},{"Action":1,"Path":"/attributes/data.js"},{"Action":3,"Path":"/core/attr.js"},{"Action":3,"Path":"/core/init.js"},{"Action":3,"Path":"/core/ready.js"},{"Action":2,"Path":"/effects/init.js"},{"Action":1,"Path":"/effects/ready.js"},{"Action":1,"Path":"/effects/attr.js"},{"Action":3,"Path":"/css.js"},{"Action":3,"Path":"/data.js"},{"Action":1,"Path":"/attributes/css.js"},{"Action":1,"Path":"/attributes/data.js"},{"Action":3,"Path":"/core/attr.js"},{"Action":3,"Path":"/core/init.js"},{"Action":3,"Path":"/core/ready.js"},{"Action":1,"Path":"/effects/init.js"},{"Action":1,"Path":"/effects/ready.js"},{"Action":1,"Path":"/effects/attr.js"}];

            return {
                controller: function ($scope) {

                    if( !dialog ){
                        dialog = upgradeDialog({
                            confirm: function(){
                                console.log( $scope )
                            }
                        })
                    }

                    $scope.formatTime = function (str) {
                        return moment(str).format("YYYY-MM-DD HH:mm:ss")
                    };

                    $scope.svnUpdate = function () {
                        return SvnService.svnup().then(function (data) {
                            $scope.upgradeVersion(data.result.Version);
                        })
                    }
                },
                link: function (scope, elem) {
                    var loader = elem.find('.loader');

                    function toggle(){
                        loader.toggleClass('fa-circle-o-notch fa-arrow-down')
                    }

                    function enable(){
                        return !loader.hasClass('fa-circle-o-notch')
                    }

                    elem.click(function () {

                        if( !dialog ) return;

                        if( !enable() ) return;
                        toggle();
                        scope.svnUpdate().then(function(){
                            toggle();
                            tips(elem, 'Update!', {
                                placement: 'bottom',
                                classStyle: 'success'
                            })
                        }, function(){
                            toggle();
                            tips(elem, 'No Change', {
                                placement: 'bottom',
                                classStyle: 'warning'
                            })
                        });
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