define([
'kernel',
'angular',
'./module',
'components/form/FormFlyout',
'components/form/FormDialog',
'ui/confirm',
'ui/tips',
'service/GroupService',
'service/ClientService'
],
function( core, ng, directive, FormFlyout, FormDialog, confirm, tips){

    directive
        .factory('EditClientDialog', function(ClientService){
            var formDialog = FormDialog({
                useLabel: true,
                buttons: ['append', {
                    text: 'Delete',
                    className: 'btn-danger',
                    click: function ( comp, e ) {
                        var anchor = e.target;
                        confirm(anchor, function(){
                            info(anchor, 'upload complete')
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
                    }else {
                        this.setTitle('Edit Client');
                    }
                },
                onShown: function(){
                    this.focusFirst();
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
                                    return;
                                }
                            });
                            same ? self.hide() : ClientService.update(data.id, data).then(function(data){
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
                controller: function( $scope ){
                    $scope.clientSelectable = false;

                    $scope.setClientSelectable = function( enable ){
                        $scope.clientSelectable = enable;
                        $scope.$digest();
                    };

                    $scope.getSelectedClient = function(){
                        $scope.mapClients(function( client ){

                        })
                    }();
                }
            }
        })
        .directive('clientSelect', function(){
            return {
                link: function( scope, elem ){
                    var checkbox = elem.find('input');
                    checkbox.change(function(){
                        scope.client._selected = checkbox.is(':checked');
                    });
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
                            .show('add');
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
        });
});