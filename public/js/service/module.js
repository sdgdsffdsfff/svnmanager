define(['angular'], function(ng){
    return ng.module('App.services',[])
        .factory('Helper', function($q){
            return {
                result : function( http ){
                    var service = this;
                    var q = $q.defer();

                    http.then(function(data){
                        service.resultData(data, q)
                    }, function(err){
                        q.reject(err)
                    });

                    return q.promise;
                },
                resultData: function(data, q) {
                    if( data && data.status == 200 ){
                        data = data.data
                    }else {
                        q.reject(data)
                    }

                    if( data.code == 'success' ){
                        q.resolve(data)
                    }else{
                        q.reject(data)
                    }
                }
            }
        })
});