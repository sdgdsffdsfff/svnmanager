define([
'kernel',
'angular',
'./module',
'moment',
'ui/Dialog',
'ui/tips',
'ui/confirm',
'ui/Toast',
'react',
'components/ui/upgradeDialog',
'service/MasterService',
'service/GlobalControlUI'
],
function( core, ng, directive, moment, Dialog, tips, confirm, Toast, React, upgradeDialog ){
    directive
        .factory('DeployDialog', function( MasterService ){
            var dialog = upgradeDialog({
                confirm: function( btn ){
                    btn.loading();
                    this.check().then(function(){
                        btn.reset();
                    }, function( data ){
                        btn.reset();
                        tips(btn.$elem(), data.message, 'warning');
                    });
                }
            }, null, {
                check: function(){
                    var self = this;
                    return this.upfileList.message().then(function( msg ){
                        return self.upfileList.getReadyToDeployFile().then(function( list ){
                            self.hide();
                            core.delay(function(){
                                self.scope.readyToDeploy( msg, list );
                            }, 500)
                        })
                    });
                },
                notify: function( text ){
                    this.upfileList.notify(text);
                },
                getUnDeployFiles: function(){
                    var self = this;
                    return MasterService.getUndeployFileList().then(function( data ){
                        self.upfileList.setList(data.result);
                        self.show();
                    })
                },
                scope: null,
                setScope: function( s ){
                    this.scope = s;
                }
            });

            return dialog;
        })
        .directive('update', function (MasterService, DeployDialog) {
            return {
                controller: function ($scope) {

                    $scope.formatTime = function (str) {
                        return moment(str).format("YYYY-MM-DD HH:mm:ss")
                    };

                    $scope.svnUpdate = function () {
                        return MasterService.svnup().then(function (data) {
                            DeployDialog.getUnDeployFiles();
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

                        scope.svnUpdate().then(function( data ){
                            toggle();
                            tips(elem, 'Update To New Version!', {
                                placement: 'bottom',
                                classStyle: 'success'
                            })
                        }, function(){
                            toggle();
                            tips(elem, 'No changes Detected!', {
                                placement: 'bottom',
                                classStyle: 'warning'
                            })
                        });
                    });
                }
            }
        })
        .directive('compile', function( MasterService ){
            return {
                link: function( scope, elem ){
                    elem.click(function(){
                        MasterService.compile();
                    })
                }
            }
        })
        .directive('deploy', function (ClientService, MasterService, DeployDialog, GlobalControlUI) {
            return {
                controller: function( $scope ){
                    $scope.readyToDeploy = function( msg, filesId ){
                        $scope.clientSelectable = true;

                        //选择在线的机
                        GlobalControlUI.show('Select the client which you want to deploy.', function(){

                            var ids = $scope.getSelectedClient().map(function( client ){
                                return client.Id
                            });

                            if( ids.length ){
                                if( filesId.length > 0 && filesId[0] !== 0 ) {
                                    filesId = filesId.map(function( t ){
                                        return t.Id;
                                    });
                                }
                                MasterService.deploy(filesId, ids, msg).then(function (data) {
                                    GlobalControlUI.hide();
                                    $scope.clientSelectable = false;
                                    //$scope.setClientSelectable(false);
                                    $.each(data.result, function( id, value ){
                                        $scope.findClient(id).Version = value.Version;
                                    })
                                })
                            }else{
                                tips(GlobalControlUI.$nextBtn, 'No client selected!', 'warning');
                            }
                        }, function(){
                            $scope.clientSelectable = false;
                            GlobalControlUI.hide();
                        });
                    };

                    DeployDialog.setScope( $scope );
                },
                link: function (scope, elem) {
                    elem.click(function () {
                        DeployDialog.getUnDeployFiles().then($.noop, function(){
                            tips(elem, 'Noting to commit!', {
                                placement: 'bottom',
                                classStyle: 'warning'
                            })
                        });
                    })
                }
            }
        })
});