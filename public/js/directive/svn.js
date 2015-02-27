define([
'kernel',
'angular',
'./module',
'moment',
'ui/Dialog',
'ui/tips',
'ui/confirm',
'react',
'components/ui/upgradeDialog',
'service/SvnService'
],
function( core, ng, directive, moment, Dialog, tips, confirm, React, upgradeDialog ){
    directive
        .factory('DeployDialog', function( SvnService ){
            var dialog = upgradeDialog({
                confirm: function( btn ){
                    confirm(btn.$elem(), function(){
                        console.log(11)
                    });
                }
            }, null, {
                deploy: function (ids) {
                    return SvnService.deploy(ids).then(function (data) {
                        console.log(data)
                    }, function(data){
                        console.log(data)
                    })
                },
                getUndeployFiles: function(){
                    var self = this;
                    return SvnService.getUndeployFileList().then(function( data ){
                        self.upfileList.setList(data.result);
                        self.show();
                    })
                }
            });

            return dialog;
        })
        .directive('svnUpdate', function (SvnService, DeployDialog) {
            return {
                controller: function ($scope) {

                    $scope.formatTime = function (str) {
                        return moment(str).format("YYYY-MM-DD HH:mm:ss")
                    };

                    $scope.svnUpdate = function () {
                        return SvnService.svnup().then(function (data) {
                            var version = data.result;
                            $scope.upgradeVersion({
                                Version: version.Version,
                                Time: version.Time
                            });
                            DeployDialog.getUndeployFiles();
                            return data;
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

                        if( !enable() ) return;
                        toggle();

                        scope.svnUpdate().then(function(){
                            toggle();
                            tips(elem, 'Update To New Version!', {
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
        .directive('svnDeploy', function (SvnService, DeployDialog) {
            return {
                link: function (scope, elem) {
                    elem.click(function () {
                        DeployDialog.getUndeployFiles().then(null, function(){
                            tips(elem, 'No files need to be deploy', {
                                placement: 'bottom',
                                classStyle: 'warning'
                            })
                        })
                    })
                }
            }
        })
});