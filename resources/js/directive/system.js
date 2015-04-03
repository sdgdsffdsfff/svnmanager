/**
 * Created by languid on 3/12/15.
 */

define([
'kernel',
'angular',
'./module',
'components/form/FormDialog'
],
function( core, ng, directive, FormDialog ){
    directive
        .directive('systemSetting', function(){
            var dialog = FormDialog({
                useLabel: true,
                fields: [{
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
                    name: 'SourceUrl',
                    placeholder: 'source url'
                }, {

                }]
            });

            return {
                controller: function( $scope ){

                },
                link: function( scope, elem ){
                    elem.click(function(){
                        dialog.show();
                    })
                }
            }
        });
});

