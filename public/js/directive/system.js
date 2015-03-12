/**
 * Created by languid on 3/12/15.
 */

define([
'kernel',
'angular',
'./module',
],
function( core, ng, directive ){
    directories
        .directive('systemSetting', function(){
            return {
                controller: function( $scope ){

                },
                link: function( scope, elem ){

                }
            }
        });
});

