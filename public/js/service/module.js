define([
'angular',
'ngWebSocket'
],
function(ng){
    return ng.module('App.services',['ngWebSocket'])
        .factory('Helper', function($q){
            return {
                result : function( http ){
                    var service = this;
                    var q = $.Deferred();

                    http.then(function(data){
                        service.resultData(data, q)
                    }, function(err){
                        q.reject(err)
                    });

                    return q;
                },
                resultData: function(data, q) {
                    if( data && data.status == 200 ){
                        data = data.data;
                    }else {
                        q.reject(data)
                    }
                    this.data(data, q)
                },
                data: function(data , q){
                    if( ng.isUndefined(q) ){
                        q = $.Deferred();
                    }

                    if( data.code == 'success' ){
                        q.resolve(data)
                    }else{
                        q.reject(data)
                    }

                    return q;
                }
            }
        })
});