/**
 * Created by languid on 11/26/14.
 */

define(['kernel', 'angular', './module'],
function( core, ng, service){
    service.factory('SvnService', ['$http', 'Helper',function( $http, Helper ){
        return {
            svnup: function () {
                return Helper.result(
                    $http.post("/aj/svn/up", {
                        paths: "_res,inc,ruochu/mobile"
                    })
                )
            },
            getLastVersion: function(){
                return Helper.result(
                    $http.get('/aj/svn/lastVersion')
                );
            },
            deploy: function( files ){
                return Helper.result(
                    $http.post("/aj/svn/deploy", {
                        files: files || []
                    })
                )
            }
        }
    }])
});