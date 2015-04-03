define([
'kernel',
'angular',
'./module',
'components/form/FormFlyout',
'components/form/FormDialog',
'components/ui/revertDialog',
'ui/Dialog',
'ui/confirm',
'ui/tips',
'ui/Flyout',
'service/GroupService',
'service/ClientService'
],
function( core, ng, directive, FormFlyout, FormDialog, revertDialog, Dialog, confirm, tips, Flyout){

    directive
        .factory('EditClientDialog', function(ClientService){
            var formDialog = FormDialog({
                useLabel: true,
                buttons: ['append', {
                    text: 'Delete',
                    name: 'Delete',
                    className: 'btn-danger',
                    loadingText: 'Deleting',
                    click: function ( btn ) {
                        var self = this;
                        var scope = this.getScope(), id = scope.client.Id;
                        confirm(btn.$elem(), function(){
                            btn.loading();
                            ClientService.del(id).then(function(){
                                scope.delClient(id);
                                self.hide();
                                btn.reset();
                            }, function(){
                                btn.reset();
                            });
                        })
                    }
                }],
                fields: [{
                    name: 'Name',
                    placeholder: 'client name'
                }, {
                    name: 'Ip',
                    required: true,
                    placeholder: 'ip address'
                }, {
                    name: 'InternalIp',
                    placeholder: 'intranet ip',
                    require: true
                }, {
                    name: 'Port',
                    placeholder: 'port'
                }, {
                    name: 'DeployPath',
                    placeholder: 'deploy path'
                }, {
                    name: 'BackupPath',
                    placeholder: 'backup path'
                },{
                    name: 'Id',
                    type: 'hidden'
                }]
            }, {
                title: 'Add Host',
                classStyle: 'medium',
                onShow: function () {
                    if( this.getState() == 'add' ){
                        this.clearForm();
                        this.setTitle('Add Client');
                        this.getRef('formBtns').getRef('Delete').$elem().hide();
                    }else {
                        this.setTitle('Edit Client');
                        this.getRef('formBtns').getRef('Delete').$elem().show();
                    }
                },
                onHidden: function(){
                    this.setState(null);
                }
            }, {
                _state: null,
                _scope: null,
                getState: function(){
                    return this._state;
                },
                setState: function( state ){
                    this._state = state;
                    return this;
                },
                setScope: function( scope ){
                    this._scope = scope;
                    return this;
                },
                getScope: function(){
                    return this._scope;
                },
                submitForm: function () {
                    var self = this;
                    var scope = this.getScope();
                    this.getFormData().then(function (data) {

                        if( self.getState() == 'add' ){
                            delete data.Id;
                            data['group'] = scope.group.Id;
                            ClientService.add(data).then(function(data){
                                scope.addClientToGroup(data.result, data.result.Group);
                                self.hide();
                            }, function(data){
                                console.log( data )
                            })
                        }else{
                            data.Id *= 1;
                            var same = true;
                            $.each(data, function(key, value){
                                if( scope.client[key] != value ){
                                    same = false;
                                    return false;
                                }
                            });
                            same ? self.hide() : ClientService.edit(data.Id, data).then(function(data){
                                //直接更新对象引用
                                scope.client = data.result;
                                self.hide();
                            })

                        }
                    });
                    return this;
                }
            });
            return formDialog;
        })
        .directive('groupAdd', function (GroupService) {

            var formFlyout = FormFlyout({
                title: 'Add a group',
                useLabel: true,
                fields: [{
                    name: 'name',
                    placeholder: 'Group Name',
                    required: true,
                    keyup: function (input, e) {
                        if (e.keyCode == 13) {
                            this.submitForm();
                        }
                    }
                }]
            }, {
                onShow: function () {
                    this.clearForm();
                    this.focusFirst();
                }
            }, {
                submitForm: function () {
                    var flyout = this;

                    this.getFormData().then(function (data) {
                        GroupService.add(data.name).then(function (data) {
                            //TODO add to group list
                            flyout.hide();
                        }, function (error) {

                        });
                    })
                }
            });

            return {
                link: function (scope, elem) {
                    elem.click(function () {
                        formFlyout.show(this, 'bottom', 'right')
                    })
                }
            }
        })
        .directive('groupList', function(){
            return {
                controller: function( $scope, $timeout ){
                    $scope.clientSelectable = false;

                    $scope.getSelectedClient = function(){
                        var clients = [];
                        $scope.mapClients(function( client ){
                            if( client.selected ){
                                clients.push(client)
                            }
                        });
                        return clients;
                    };

                    $scope.notify = function(id, msg, hide){
                        var client = $scope.findClient(id);
                        if( client ) {
                            client.Message = msg;
                            if( hide ){
                                $timeout(function(){
                                    client.Message = null;
                                }, 2000)
                            }
                        }
                    }
                }
            }
        })
        .factory('ClientControlUI', function(){
            var lastBindName, lastElem = null;
            return {
                show: function( elem ){
                    this.hide();
                    lastElem = elem.addClass('open');
                    lastBindName = core.clickAnyWhereHideButMe(elem, function(){
                        elem.removeClass('open');
                        lastElem = null;
                    });
                },
                hide: function(){
                    if( lastElem && lastElem.hasClass('open') ){
                        lastElem.removeClass('open');
                    }
                    if( lastBindName ){
                        core.unbindDocumentEvent(lastBindName);
                    }
                }
            }
        })
        .directive('clientControl', function( ClientControlUI ){
            return {
                link: function( scope, elem ){
                    var more = elem.find('em.more');
                    more.on('click',function(){
                        ClientControlUI.show(elem)
                    })
                }
            }
        })
        .directive('clientMove', function( ClientService ){
            var formFlyout = FormFlyout({
                title: "Change client group",
                fields: [{
                    name: 'Group',
                    type: 'select'
                }]
            }, null, {
                watched: false,
                _scope: null,
                setScope: function( scope ){
                    this._scope = scope;
                    return this;
                },
                getScope: function(){
                    return this._scope;
                },
                selectGroup: function( id ){
                    this.getRef('formBody').getRef('Group').setValue(id);
                    return this;
                },
                submitForm: function () {
                    var self = this;
                    this.getFormData().then(function (data) {
                        var scope = self.getScope();
                        var gid = scope.$parent.group.Id,
                            to = data.Group * 1;

                        if (gid == to ) return;

                        scope.$parent.moveClientToGroup(scope.client.Id, gid, to);

                        ClientService.changeGroup(scope.client.Id, to).then(function () {
                            self.hide();
                        });
                    });

                    return this;
                },
                $scope: null
            });

            return {
                controller: function( $scope ){
                    if( !formFlyout.watched ) {
                        formFlyout.watched = true;
                        $scope.onGroupChange(function( values ){
                            formFlyout.getRef('formBody').getRef('Group').setValues(values);
                        })
                    }
                },
                link: function( scope, elem ){
                    elem.click(function(){
                        formFlyout
                            .setScope(scope)
                            .selectGroup( scope.group.Id )
                            .show(this, 'bottom', 'right');
                    })
                }
            }
        })
        .directive('clientAdd', function ( EditClientDialog ) {
            return {
                link: function (scope, elem) {
                    elem.click(function () {
                        EditClientDialog
                            .setScope(scope)
                            .setState('add')
                            .show();
                    });
                }
            }
        })
        .directive('clientEdit', function (EditClientDialog) {

            return {
                scope: '=',
                link: function (scope, elem) {
                    elem.click(function () {
                        EditClientDialog
                            .setScope(scope)
                            .setState('edit')
                            .setFormValue(scope.client)
                            .show();
                    })
                }
            }
        })
        .directive('clientUpdate', function( ClientService ){
            return {
                link: function( scope, elem ){
                    elem.click(function(){
                        ClientService.update( scope.client.Id ).then(function( data ){
                            scope.client.Version = data.result.Version;
                        }, function( data ){
                            console.log( data )
                        })
                    })
                }
            }
        })
        .directive('clientDeploy', function( ClientService, ClientControlUI) {
            var lastDefer;

            return {
                link: function( scope, elem ){
                    elem.click(function(){
                        if( lastDefer && lastDefer.state() == 'pending' ){
                            tips(elem, 'is busy', 'info');
                            return
                        }
                        lastDefer = ClientService.deploy( scope.client.Id ).then(function( data ){
                            ClientControlUI.hide();
                        }, function(){
                            tips(elem, 'is busy', 'info')
                        })
                    })
                }
            }
        })
        .directive('clientRevert', function( ClientService ) {
            var lastDefer;
            var dialog = revertDialog({
                revert: function( path ){
                    ClientService.revert(this.scope.client.Id, path).then(function( data ){
                        console.log( data )
                    })
                },
                remove: function( path ){
                    ClientService.removeBackup(this.scope.client.Id, path).then(function( data ){
                        console.log( data )
                    })
                }
            }, null, {
                scope: null,
                setScope: function( scope ){
                    this.scope = scope;
                    return this
                }
            });

            return {

                link: function( scope, elem ){
                    elem.click(function(){
                        if( lastDefer && lastDefer.state() == 'pending' ){
                            tips(elem, 'fetching..', 'info');
                            return
                        }
                        lastDefer = ClientService.getBackupList( scope.client.Id ).then(function( data ){
                            var list = data.result;
                            list.sort().reverse();
                            dialog.revertList.setList(list);
                            dialog.setScope(scope).show();
                        }, function( data ){
                            tips(elem, 'error', 'warning')
                        })
                    });
                }
            }
        })
        .directive('clientError', function(){
            return {
                link: function( scope, elem ){
                    elem.click(function(){
                        tips(elem, scope.client.Error, 'danger')
                    })
                }
            }
        })
        .directive('clientLog', function( ClientService, LogDialogUI ) {
            var lastDefer;

            return {
                link: function( scope, elem ){
                    elem.click(function(){
                        if( lastDefer && lastDefer.state() == 'pending' ){
                            tips(elem, 'waiting', 'info');
                            return;
                        }
                        lastDefer = ClientService.log( scope.client.Id ).then(function( data ){
                            LogDialogUI.setTitle('Client log').setContent(data.message).show();
                        }, function( data ){
                            console.error('client:'+scope.client.Id+'/'+scope.client.Ip, data.message);
                            tips(elem, 'Can not open log file!', 'warning');
                        });
                    })
                }
            }
        })
});